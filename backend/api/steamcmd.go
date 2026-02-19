package api

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

	// Windows 隐藏控制台窗口
	hideWindow(cmd)

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
	var mkdirErr error
	if runtime.GOOS == "windows" {
		mkdirErr = os.MkdirAll(cfg.SteamCMDPath, 0755)
	} else {
		mkdirErr = os.MkdirAll(cfg.SteamCMDPath, 0755)
	}
	if mkdirErr != nil {
		fail(c, "创建目录失败: "+mkdirErr.Error())
		return
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		ws.Broadcast("开始下载 SteamCMD (Windows)...")
		// 使用 -ExecutionPolicy Bypass 绕过策略限制
		// 使用 -UseBasicParsing 兼容没有 IE 的环境
		script := `$ProgressPreference = 'SilentlyContinue'
$url = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
$zip = "$env:TEMP\steamcmd.zip"
Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
Expand-Archive -Path $zip -DestinationPath "` + cfg.SteamCMDPath + `" -Force
Remove-Item $zip`
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
		hideWindow(cmd)
	} else {
		ws.Broadcast("开始下载 SteamCMD (Linux)...")
		script := `cd "` + cfg.SteamCMDPath + `" && \
curl -sqL "https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz" | tar zxvf - && \
chmod +x steamcmd.sh`
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
		hideWindow(cmd)
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
	status.Running = isProcessRunning()

	success(c, status)
}

// isProcessRunning 检查 Arma Reforger 服务端是否运行
func isProcessRunning() bool {
	if runtime.GOOS == "windows" {
		// Windows: 使用 tasklist 查找进程
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq ArmaReforgerServer.exe", "/FO", "CSV", "/NH")
		hideWindow(cmd)
		out, err := cmd.Output()
		if err != nil {
			return false
		}
		return strings.Contains(string(out), "ArmaReforgerServer.exe")
	} else {
		// Linux: 使用 pgrep 查找进程
		cmd := exec.Command("pgrep", "-f", "ArmaReforgerServer")
		out, err := cmd.Output()
		return err == nil && len(strings.TrimSpace(string(out))) > 0
	}
}

// getProcessPID 获取 Arma Reforger 服务端 PID
func getProcessPID() int {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq ArmaReforgerServer.exe", "/FO", "CSV", "/NH")
		hideWindow(cmd)
		out, err := cmd.Output()
		if err != nil {
			return 0
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.Contains(line, "ArmaReforgerServer.exe") {
				// CSV 格式: "ArmaReforgerServer.exe","1234","..."
				parts := strings.Split(line, "\",\"")
				if len(parts) >= 2 {
					pidStr := strings.Trim(parts[1], "\"")
					pid, _ := strconv.Atoi(pidStr)
					return pid
				}
			}
		}
		return 0
	} else {
		cmd := exec.Command("pgrep", "-f", "ArmaReforgerServer")
		out, err := cmd.Output()
		if err != nil {
			return 0
		}
		pidStr := strings.TrimSpace(string(out))
		pid, _ := strconv.Atoi(pidStr)
		return pid
	}
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

	cmd := exec.Command(
		steamcmd,
		"+force_install_dir", cfg.ServerPath,
		"+login", "anonymous",
		"+app_update", "1874900", "validate",
		"+quit",
	)

	if runtime.GOOS == "windows" {
		hideWindow(cmd)
	}

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
