package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UserManager 用户管理器
type UserManager struct {
	mu       sync.RWMutex
	dataDir  string
	users    map[string]*UserInfo
	enabled  bool
}

// UserInfo 用户详细信息
type UserInfo struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"` // admin, user
	CreatedAt    int64  `json:"created_at"`
	LastLoginAt  int64  `json:"last_login_at,omitempty"`
}

// UsersConfig 用户配置文件结构
type UsersConfig struct {
	Enabled  bool                 `json:"enabled"`
	Users    map[string]*UserInfo `json:"users"`
}

var (
	manager *UserManager
	once    sync.Once
)

// GetUserManager 获取用户管理器单例
func GetUserManager() *UserManager {
	once.Do(func() {
		// 默认数据目录
		dataDir := "./data"
		if envDir := os.Getenv("ARSM_DATA_DIR"); envDir != "" {
			dataDir = envDir
		}
		manager = NewUserManager(dataDir)
		
		// 加载配置，失败时输出错误但继续运行（降级到认证禁用状态）
		if err := manager.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "[Auth] ❌ 加载用户配置失败: %v\n", err)
			fmt.Fprintf(os.Stderr, "[Auth] ⚠️ 认证功能将被禁用\n")
			manager.enabled = false
		}
	})
	return manager
}

// NewUserManager 创建用户管理器
func NewUserManager(dataDir string) *UserManager {
	return &UserManager{
		dataDir: dataDir,
		users:   make(map[string]*UserInfo),
		enabled: false,
	}
}

// getConfigPath 获取配置文件路径
func (um *UserManager) getConfigPath() string {
	return filepath.Join(um.dataDir, "users.json")
}

// Load 加载用户配置
func (um *UserManager) Load() error {
	um.mu.Lock()
	defer um.mu.Unlock()

	configPath := um.getConfigPath()
	
	// 确保目录存在
	if err := os.MkdirAll(um.dataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %w", err)
	}

	// 如果文件不存在，创建默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("[Auth] 配置文件不存在，正在创建默认配置: %s\n", configPath)
		
		// 创建默认管理员账户
		hash, err := HashPassword("admin")
		if err != nil {
			return fmt.Errorf("生成默认密码失败: %w", err)
		}
		
		um.users["admin"] = &UserInfo{
			Username:     "admin",
			PasswordHash: hash,
			Role:         "admin",
			CreatedAt:    time.Now().Unix(),
		}
		um.enabled = true
		
		if err := um.saveLocked(); err != nil {
			return fmt.Errorf("保存默认配置失败: %w", err)
		}
		
		fmt.Printf("[Auth] ✅ 默认配置创建成功 (用户名: admin, 密码: admin)\n")
		return nil
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config UsersConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	um.enabled = config.Enabled
	um.users = config.Users
	
	// 如果没有用户，创建默认管理员
	if len(um.users) == 0 {
		hash, _ := HashPassword("admin")
		um.users["admin"] = &UserInfo{
			Username:     "admin",
			PasswordHash: hash,
			Role:         "admin",
			CreatedAt:    time.Now().Unix(),
		}
		um.enabled = true
		return um.saveLocked()
	}

	return nil
}

// Save 保存用户配置
func (um *UserManager) Save() error {
	um.mu.Lock()
	defer um.mu.Unlock()
	return um.saveLocked()
}

// saveLocked 保存（内部使用，已加锁）
func (um *UserManager) saveLocked() error {
	config := UsersConfig{
		Enabled: um.enabled,
		Users:   um.users,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	configPath := um.getConfigPath()
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败 %s: %w", configPath, err)
	}

	return nil
}

// IsEnabled 检查认证是否启用
func (um *UserManager) IsEnabled() bool {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return um.enabled
}

// SetEnabled 设置认证启用状态
func (um *UserManager) SetEnabled(enabled bool) error {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.enabled = enabled
	return um.saveLocked()
}

// Authenticate 验证用户登录
func (um *UserManager) Authenticate(username, password string) (*UserInfo, error) {
	um.mu.RLock()
	defer um.mu.RUnlock()

	// 如果认证未启用，直接通过
	if !um.enabled {
		return nil, nil
	}

	user, exists := um.users[username]
	if !exists {
		return nil, errors.New("用户名或密码错误")
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

// GetUser 获取用户信息
func (um *UserManager) GetUser(username string) (*UserInfo, bool) {
	um.mu.RLock()
	defer um.mu.RUnlock()
	user, exists := um.users[username]
	return user, exists
}

// CreateUser 创建用户
func (um *UserManager) CreateUser(username, password, role string) error {
	if username == "" || password == "" {
		return errors.New("用户名和密码不能为空")
	}

	if role != "admin" && role != "user" {
		role = "user"
	}

	um.mu.Lock()
	defer um.mu.Unlock()

	if _, exists := um.users[username]; exists {
		return errors.New("用户已存在")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	um.users[username] = &UserInfo{
		Username:     username,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    time.Now().Unix(),
	}

	return um.saveLocked()
}

// UpdateUser 更新用户信息
func (um *UserManager) UpdateUser(username, newPassword, newRole string) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	user, exists := um.users[username]
	if !exists {
		return errors.New("用户不存在")
	}

	if newPassword != "" {
		hash, err := HashPassword(newPassword)
		if err != nil {
			return err
		}
		user.PasswordHash = hash
	}

	if newRole != "" {
		if newRole != "admin" && newRole != "user" {
			return errors.New("无效的角色")
		}
		user.Role = newRole
	}

	return um.saveLocked()
}

// DeleteUser 删除用户
func (um *UserManager) DeleteUser(username string) error {
	um.mu.Lock()
	defer um.mu.Unlock()

	if username == "admin" {
		return errors.New("不能删除默认管理员账户")
	}

	if _, exists := um.users[username]; !exists {
		return errors.New("用户不存在")
	}

	delete(um.users, username)
	return um.saveLocked()
}

// ListUsers 获取所有用户列表
func (um *UserManager) ListUsers() []*UserInfo {
	um.mu.RLock()
	defer um.mu.RUnlock()

	users := make([]*UserInfo, 0, len(um.users))
	for _, user := range um.users {
		// 返回副本，不暴露密码
		users = append(users, &UserInfo{
			Username:    user.Username,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			LastLoginAt: user.LastLoginAt,
		})
	}
	return users
}

// UpdateLastLogin 更新最后登录时间
func (um *UserManager) UpdateLastLogin(username string) {
	um.mu.Lock()
	defer um.mu.Unlock()

	if user, exists := um.users[username]; exists {
		user.LastLoginAt = time.Now().Unix()
		um.saveLocked()
	}
}

// CheckDefaultPassword 检查是否还在使用默认密码
func (um *UserManager) CheckDefaultPassword() bool {
	um.mu.RLock()
	defer um.mu.RUnlock()

	if admin, exists := um.users["admin"]; exists {
		return CheckPassword("admin", admin.PasswordHash)
	}
	return false
}
