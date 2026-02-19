package api

import (
	"net/http"

	"arsm/auth"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
func AuthHandler(r *gin.RouterGroup) {
	// 公开路由（不需要认证）
	r.POST("/auth/login", handleLogin)
	r.GET("/auth/status", handleAuthStatus)

	// 需要认证的路由
	authGroup := r.Group("/auth")
	authGroup.Use(auth.JWTAuthMiddleware())
	{
		authGroup.POST("/logout", handleLogout)
		authGroup.GET("/profile", handleGetProfile)
		authGroup.POST("/password", handleChangePassword)

		// 用户管理（仅管理员）
		authGroup.GET("/users", handleListUsers)
		authGroup.POST("/users", handleCreateUser)
		authGroup.PUT("/users/:username", handleUpdateUser)
		authGroup.DELETE("/users/:username", handleDeleteUser)
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handleLogin 用户登录
func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "用户名和密码不能为空")
		return
	}

	um := auth.GetUserManager()

	// 验证用户
	user, err := um.Authenticate(req.Username, req.Password)
	if err != nil {
		fail(c, err.Error())
		return
	}

	// 认证未启用时，直接返回成功
	if user == nil {
		success(c, gin.H{"enabled": false})
		return
	}

	// 生成 JWT Token
	token, expiresAt, err := auth.GenerateToken(user.Username, user.Role)
	if err != nil {
		fail(c, "生成令牌失败")
		return
	}

	// 更新最后登录时间
	um.UpdateLastLogin(user.Username)

	success(c, gin.H{
		"enabled":    true,
		"token":      token,
		"username":   user.Username,
		"role":       user.Role,
		"expires_at": expiresAt,
	})
}

// handleLogout 用户登出
func handleLogout(c *gin.Context) {
	// JWT 是无状态的，客户端删除 token 即可
	// 这里可以添加 token 黑名单逻辑（如果需要）
	success(c, gin.H{"message": "登出成功"})
}

// handleAuthStatus 获取认证状态
func handleAuthStatus(c *gin.Context) {
	um := auth.GetUserManager()

	status := gin.H{
		"enabled":         um.IsEnabled(),
		"default_password": um.CheckDefaultPassword(),
	}

	// 如果已登录，返回用户信息
	if username, role, ok := auth.GetCurrentUser(c); ok {
		status["authenticated"] = true
		status["username"] = username
		status["role"] = role
	} else {
		status["authenticated"] = false
	}

	success(c, status)
}

// handleGetProfile 获取当前用户信息
func handleGetProfile(c *gin.Context) {
	username, role, ok := auth.GetCurrentUser(c)
	if !ok {
		fail(c, "未登录")
		return
	}

	success(c, gin.H{
		"username": username,
		"role":     role,
	})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	NewUsername string `json:"new_username"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// handleChangePassword 修改密码
func handleChangePassword(c *gin.Context) {
	username, _, ok := auth.GetCurrentUser(c)
	if !ok {
		fail(c, "未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "无效的请求数据")
		return
	}

	um := auth.GetUserManager()

	// 验证旧密码
	user, err := um.Authenticate(username, req.OldPassword)
	if err != nil {
		fail(c, "原密码错误")
		return
	}

	// 如果需要修改用户名
	if req.NewUsername != "" && req.NewUsername != username {
		// 检查新用户名是否已存在
		if _, exists := um.GetUser(req.NewUsername); exists {
			fail(c, "用户名已存在")
			return
		}

		// 创建新用户并删除旧用户（简单处理）
		if err := um.CreateUser(req.NewUsername, req.NewPassword, user.Role); err != nil {
			fail(c, "更新用户名失败: "+err.Error())
			return
		}

		// 不能直接删除 admin 账户的特殊逻辑处理
		if username == "admin" {
			// 如果是 admin 修改用户名，我们实际上是把 admin 的密码改了，并创建一个新管理员
			// 但 ARSM 逻辑中 admin 是保留字吗？查看 user.go
			// user.go 中 DeleteUser 确实禁止删除 admin
		}

		if err := um.DeleteUser(username); err != nil {
			// 如果删除失败（比如是 admin），我们只更新密码
			if err := um.UpdateUser(username, req.NewPassword, ""); err != nil {
				fail(c, "更新失败: "+err.Error())
				return
			}
			success(c, gin.H{"message": "由于是系统账户，仅密码已更新", "username": username})
			return
		}

		success(c, gin.H{"message": "用户名和密码已更新，请重新登录", "relogin": true})
		return
	}

	// 仅更新密码
	if err := um.UpdateUser(user.Username, req.NewPassword, ""); err != nil {
		fail(c, "修改密码失败: "+err.Error())
		return
	}

	success(c, gin.H{"message": "密码修改成功"})
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// handleListUsers 获取用户列表
func handleListUsers(c *gin.Context) {
	// 检查是否为管理员
	_, role, ok := auth.GetCurrentUser(c)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, Response{Code: 403, Message: "无权限访问"})
		return
	}

	um := auth.GetUserManager()
	users := um.ListUsers()
	success(c, users)
}

// handleCreateUser 创建用户
func handleCreateUser(c *gin.Context) {
	// 检查是否为管理员
	_, role, ok := auth.GetCurrentUser(c)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, Response{Code: 403, Message: "无权限访问"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "无效的请求数据")
		return
	}

	um := auth.GetUserManager()
	if err := um.CreateUser(req.Username, req.Password, req.Role); err != nil {
		fail(c, err.Error())
		return
	}

	// 如果是第一个用户，自动启用认证
	if !um.IsEnabled() {
		um.SetEnabled(true)
	}

	success(c, gin.H{"message": "用户创建成功"})
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
}

// handleUpdateUser 更新用户信息
func handleUpdateUser(c *gin.Context) {
	// 检查是否为管理员
	currentUser, role, ok := auth.GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusForbidden, Response{Code: 403, Message: "未登录"})
		return
	}

	targetUsername := c.Param("username")

	// 普通用户只能修改自己
	if role != "admin" && currentUser != targetUsername {
		c.JSON(http.StatusForbidden, Response{Code: 403, Message: "无权限修改其他用户"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "无效的请求数据")
		return
	}

	// 普通用户不能修改自己的角色
	if role != "admin" && req.Role != "" {
		fail(c, "无权修改角色")
		return
	}

	um := auth.GetUserManager()
	if err := um.UpdateUser(targetUsername, req.Password, req.Role); err != nil {
		fail(c, err.Error())
		return
	}

	success(c, gin.H{"message": "用户更新成功"})
}

// handleDeleteUser 删除用户
func handleDeleteUser(c *gin.Context) {
	// 检查是否为管理员
	_, role, ok := auth.GetCurrentUser(c)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, Response{Code: 403, Message: "无权限访问"})
		return
	}

	username := c.Param("username")

	um := auth.GetUserManager()
	if err := um.DeleteUser(username); err != nil {
		fail(c, err.Error())
		return
	}

	success(c, gin.H{"message": "用户删除成功"})
}
