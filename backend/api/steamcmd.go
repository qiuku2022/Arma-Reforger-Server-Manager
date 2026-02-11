package api

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"arsm/config"
	"arsm/models"
	"arsm/ws"

	"github.com/gin-gonic/gin"
)

// 辅助函数：执行命令并将输出实时广播到 WebSocket
func execAndStream(cmd *exec.Cmd) error {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	// 异步读取 stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			ws.Broadcast(scanner.Text())
		}
	}()

	// 异步读取 stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			ws.Broadcast("ERROR: " + scanner.Text())
		}
	}()

	return cmd.Wait()
}

// GetSteamCMDStatus 获取SteamCMD状态
func GetSteamCMDStatus(c *gin.Context) {
	cfg := config.Get()
	status := models.SteamCMDStatus{
		Installed: false,
		Path:      cfg.SteamCMDPath,
	}

	var executable string
	if runtime.GOOS == "windows" {
		executable = filepath.Join(cfg.SteamCMDPath, "steamcmd.exe")
	} else {
		executable = filepath.Join(cfg.SteamCMDPath, "steamcmd.sh")
	}

	if _, err := os.Stat(executable); err == nil {
		status.Installed = true
	}

	success(c, status)
}

// InstallSteamCMD 安装SteamCMD
func InstallSteamCMD(c *gin.Context) {
	cfg := config.Get()

	// 创建目录
	if err := os.MkdirAll(cfg.SteamCMDPath, 0755); err != nil {
		fail(c, "创建目录失败: "+err.Error())
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		ws.Broadcast("开始下载 SteamCMD (Windows)...")
		script := `
			$url = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
			$zip = "$env:TEMP\steamcmd.zip"
			Invoke-WebRequest -Uri $url -OutFile $zip
			Expand-Archive -Path $zip -DestinationPath "` + cfg.SteamCMDPath + `" -Force
			Remove-Item $zip
		`
		cmd = exec.Command("powershell", "-Command", script)
	} else {
		ws.Broadcast("开始下载 SteamCMD (Linux)...")
		script := `
			cd "` + cfg.SteamCMDPath + `" && \
			curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf - && \
			chmod +x steamcmd.sh
		`
		cmd = exec.Command("bash", "-c", script)
	}

	if err := execAndStream(cmd); err != nil {
		fail(c, "安装失败: "+err.Error())
		return
	}

	ws.Broadcast("SteamCMD 安装完成。")
	success(c, nil)
}

// UpdateSteamCMD 更新SteamCMD
func UpdateSteamCMD(c *gin.Context) {
	cfg := config.Get()
	ws.Broadcast("开始更新 SteamCMD...")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		executable := filepath.Join(cfg.SteamCMDPath, "steamcmd.exe")
		cmd = exec.Command(executable, "+login", "anonymous", "+quit")
	} else {
		executable := filepath.Join(cfg.SteamCMDPath, "steamcmd.sh")
		cmd = exec.Command(executable, "+login", "anonymous", "+quit")
	}

	if err := execAndStream(cmd); err != nil {
		fail(c, "更新失败: "+err.Error())
		return
	}

	ws.Broadcast("SteamCMD 更新完成。")
	success(c, nil)
}

// DeleteSteamCMD 删除SteamCMD
func DeleteSteamCMD(c *gin.Context) {
	cfg := config.Get()
	ws.Broadcast("正在删除 SteamCMD...")

	if err := os.RemoveAll(cfg.SteamCMDPath); err != nil {
		fail(c, "删除失败: "+err.Error())
		return
	}

	ws.Broadcast("SteamCMD 已删除。")
	success(c, nil)
}

// GetServerStatus 获取服务端状态
func GetServerStatus(c *gin.Context) {
	cfg := config.Get()
	status := models.ServerStatus{
		Installed: false,
		Running:   false,
	}

	var executable string
	if runtime.GOOS == "windows" {
		executable = filepath.Join(cfg.ServerPath, "ArmaReforgerServer.exe")
	} else {
		executable = filepath.Join(cfg.ServerPath, "ArmaReforgerServer")
	}

	if _, err := os.Stat(executable); err == nil {
		status.Installed = true
	}

	// 检查进程是否运行
	if runtime.GOOS == "windows" {
		out, _ := exec.Command("tasklist", "/FI", "IMAGENAME eq ArmaReforgerServer.exe").Output()
		status.Running = strings.Contains(string(out), "ArmaReforgerServer.exe")
	} else {
		out, _ := exec.Command("pgrep", "-f", "ArmaReforgerServer").Output()
		status.Running = len(strings.TrimSpace(string(out))) > 0
	}

	success(c, status)
}

// InstallServer 安装游戏服务端
func InstallServer(c *gin.Context) {
	cfg := config.Get()
	ws.Broadcast("开始安装/更新 Arma Reforger 服务端 (AppID 1874900)...")

	// 创建目录
	if err := os.MkdirAll(cfg.ServerPath, 0755); err != nil {
		fail(c, "创建目录失败: "+err.Error())
		return
	}

	var steamcmd string
	if runtime.GOOS == "windows" {
		steamcmd = filepath.Join(cfg.SteamCMDPath, "steamcmd.exe")
	} else {
		steamcmd = filepath.Join(cfg.SteamCMDPath, "steamcmd.sh")
	}

	cmd := exec.Command(steamcmd,
		"+force_install_dir", cfg.ServerPath,
		"+login", "anonymous",
		"+app_update", "1874900", "validate",
		"+quit",
	)

	if err := execAndStream(cmd); err != nil {
		fail(c, "安装失败: "+err.Error())
		return
	}

	ws.Broadcast("服务端安装/更新完成。")
	success(c, nil)
}

// UpdateServer 更新服务端
func UpdateServer(c *gin.Context) {
	InstallServer(c)
}

// DeleteServer 删除服务端
func DeleteServer(c *gin.Context) {
	cfg := config.Get()
	ws.Broadcast("正在删除游戏服务端文件...")

	if err := os.RemoveAll(cfg.ServerPath); err != nil {
		fail(c, "删除失败: "+err.Error())
		return
	}

	ws.Broadcast("游戏服务端已删除。")
	success(c, nil)
}
