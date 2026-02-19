package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"arsm/api"
	"arsm/auth"
	"arsm/config"
	"arsm/ws"
	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 加载配置
	config.Load()

	// 初始化用户管理器并输出日志
	um := auth.GetUserManager()
	
	// 获取绝对路径用于日志
	dataDir := "./data"
	if envDir := os.Getenv("ARSM_DATA_DIR"); envDir != "" {
		dataDir = envDir
	}
	absPath, _ := filepath.Abs(dataDir)
	
	fmt.Printf("[ARSM] 数据目录: %s\n", absPath)
	fmt.Printf("[ARSM] 认证状态: %s\n", map[bool]string{true: "已启用", false: "已禁用"}[um.IsEnabled()])
	
	if um.CheckDefaultPassword() {
		fmt.Printf("[ARSM] ⚠️ 警告: 正在使用默认密码 (admin/admin)，请尽快修改!\n")
	}

	// 生产模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 认证相关（公开）
		api.AuthHandler(apiGroup)

		// 需要认证的路由
		authorized := apiGroup.Group("/")
		authorized.Use(auth.JWTAuthMiddleware())
		{
			// 系统信息
			authorized.GET("/system/info", api.GetSystemInfo)

		// SteamCMD 管理
		authorized.GET("/steamcmd/status", api.GetSteamCMDStatus)
		authorized.POST("/steamcmd/install", api.InstallSteamCMD)
		authorized.POST("/steamcmd/update", api.UpdateSteamCMD)
		authorized.DELETE("/steamcmd", api.DeleteSteamCMD)

		// 游戏服务端管理
		authorized.GET("/server/status", api.GetServerStatus)
		authorized.POST("/server/install", api.InstallServer)
		authorized.POST("/server/update", api.UpdateServer)
		authorized.DELETE("/server", api.DeleteServer)
		authorized.POST("/server/start", api.StartServer)
		authorized.POST("/server/stop", api.StopServer)
		authorized.POST("/server/restart", api.RestartServer)

		// 配置管理
		authorized.GET("/config", api.GetConfig)
		authorized.POST("/config", api.SaveConfig)
		authorized.GET("/config/presets", api.GetPresets)
		authorized.GET("/config/presets/:name", api.GetPresetContent)
		authorized.POST("/config/presets", api.SavePreset)
		authorized.DELETE("/config/presets/:name", api.DeletePreset)
		authorized.POST("/config/import", api.ImportConfig)
		authorized.GET("/config/export", api.ExportConfig)
		authorized.GET("/config/scenarios", api.GetScenarios)

		// 模组管理
		authorized.GET("/mods", api.GetMods)
		authorized.POST("/mods", api.AddMod)
		authorized.DELETE("/mods/:id", api.DeleteMod)
		authorized.POST("/mods/:id/enable", api.EnableMod)
		authorized.POST("/mods/:id/disable", api.DisableMod)
		authorized.GET("/mods/:id/check", api.CheckModFiles)

		// RCON
		authorized.GET("/rcon/players", api.GetPlayers)
		authorized.GET("/rcon/status", api.GetRCONStatus)
		authorized.GET("/rcon/logs", api.GetRCONLogs)
		authorized.POST("/rcon/kick/:id", api.KickPlayer)
		authorized.POST("/rcon/ban/:id", api.BanPlayer)
		authorized.POST("/rcon/command", api.SendRCONCommand)

		// 设置
		authorized.GET("/settings", api.GetSettings)
		authorized.POST("/settings", api.SaveSettings)
	}
	}

	// WebSocket 日志
	r.GET("/ws/logs", ws.HandleLogs)

	// 静态文件服务
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(staticFS))))

	// 启动服务
	r.Run(":8080")
}
