package main

import (
	"embed"
	"io/fs"
	"net/http"

	"arsm/api"
	"arsm/config"
	"arsm/ws"
	"github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 加载配置
	config.Load()

	// 生产模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 系统信息
		apiGroup.GET("/system/info", api.GetSystemInfo)

		// SteamCMD 管理
		apiGroup.GET("/steamcmd/status", api.GetSteamCMDStatus)
		apiGroup.POST("/steamcmd/install", api.InstallSteamCMD)
		apiGroup.POST("/steamcmd/update", api.UpdateSteamCMD)
		apiGroup.DELETE("/steamcmd", api.DeleteSteamCMD)

		// 游戏服务端管理
		apiGroup.GET("/server/status", api.GetServerStatus)
		apiGroup.POST("/server/install", api.InstallServer)
		apiGroup.POST("/server/update", api.UpdateServer)
		apiGroup.DELETE("/server", api.DeleteServer)
		apiGroup.POST("/server/start", api.StartServer)
		apiGroup.POST("/server/stop", api.StopServer)
		apiGroup.POST("/server/restart", api.RestartServer)

		// 配置管理
		apiGroup.GET("/config", api.GetConfig)
		apiGroup.POST("/config", api.SaveConfig)
		apiGroup.GET("/config/presets", api.GetPresets)
		apiGroup.GET("/config/presets/:name", api.GetPresetContent)
		apiGroup.POST("/config/presets", api.SavePreset)
		apiGroup.DELETE("/config/presets/:name", api.DeletePreset)
		apiGroup.POST("/config/import", api.ImportConfig)
		apiGroup.GET("/config/export", api.ExportConfig)
		apiGroup.GET("/config/scenarios", api.GetScenarios)

		// 模组管理
		apiGroup.GET("/mods", api.GetMods)
		apiGroup.POST("/mods", api.AddMod)
		apiGroup.DELETE("/mods/:id", api.DeleteMod)
		apiGroup.POST("/mods/:id/enable", api.EnableMod)
		apiGroup.POST("/mods/:id/disable", api.DisableMod)
		apiGroup.GET("/mods/:id/check", api.CheckModFiles)

		// RCON
		apiGroup.GET("/rcon/players", api.GetPlayers)
		apiGroup.GET("/rcon/status", api.GetRCONStatus)
		apiGroup.GET("/rcon/logs", api.GetRCONLogs)
		apiGroup.POST("/rcon/kick/:id", api.KickPlayer)
		apiGroup.POST("/rcon/ban/:id", api.BanPlayer)
		apiGroup.POST("/rcon/command", api.SendRCONCommand)

		// 设置
		apiGroup.GET("/settings", api.GetSettings)
		apiGroup.POST("/settings", api.SaveSettings)
	}

	// WebSocket 日志
	r.GET("/ws/logs", ws.HandleLogs)

	// 静态文件服务
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.NoRoute(gin.WrapH(http.FileServer(http.FS(staticFS))))

	// 启动服务
	r.Run(":8080")
}
