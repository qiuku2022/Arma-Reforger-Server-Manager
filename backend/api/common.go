package api

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"arsm/config"
	"arsm/models"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func fail(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{Code: 1, Message: message})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	
	// CPU 使用率 (采样 200ms)
	cpuPercent, _ := cpu.Percent(200*time.Millisecond, false)
	var cpuUsage float64
	if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// 内存信息
	vMem, _ := mem.VirtualMemory()

	// 磁盘信息 (检查游戏服务端所在磁盘)
	cfg := config.Get()
	diskPath := cfg.ServerPath
	if runtime.GOOS == "windows" {
		diskPath = filepath.VolumeName(cfg.ServerPath)
		if diskPath == "" {
			diskPath = "C:"
		}
	} else {
		diskPath = "/"
	}
	
	dUsage, _ := disk.Usage(diskPath)

	info := models.SystemInfo{
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		Hostname:    hostname,
		CPUUsage:    cpuUsage,
		MemoryTotal: vMem.Total,
		MemoryUsed:  vMem.Used,
		MemoryFree:  vMem.Available,
		DiskTotal:   dUsage.Total,
		DiskUsed:    dUsage.Used,
		DiskFree:    dUsage.Free,
	}
	success(c, info)
}

// GetSettings 获取设置
func GetSettings(c *gin.Context) {
	success(c, config.Get())
}

// SaveSettings 保存设置
func SaveSettings(c *gin.Context) {
	var cfg config.AppConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		fail(c, "无效的配置数据")
		return
	}
	if err := config.Update(&cfg); err != nil {
		fail(c, "保存配置失败")
		return
	}
	success(c, nil)
}
