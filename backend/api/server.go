package api

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	"arsm/config"
	"arsm/ws"

	"github.com/gin-gonic/gin"
)

var (
	serverProcess *exec.Cmd
	serverMu      sync.Mutex
)

// StartServer 启动服务端
func StartServer(c *gin.Context) {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverProcess != nil && serverProcess.Process != nil {
		fail(c, "服务端已在运行")
		return
	}

	cfg := config.Get()

	var executable string
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		executable = filepath.Join(cfg.ServerPath, "ArmaReforgerServer.exe")
		configPath := filepath.Join(cfg.ServerPath, "config.json")
		profilePath := filepath.Join(cfg.ServerPath, "profile")
		cmd = exec.Command(executable, "-config", configPath, "-profile", profilePath)
	} else {
		executable = filepath.Join(cfg.ServerPath, "ArmaReforgerServer")
		configPath := filepath.Join(cfg.ServerPath, "config.json")
		profilePath := filepath.Join(cfg.ServerPath, "profile")
		cmd = exec.Command(executable, "-config", configPath, "-profile", profilePath)
	}

	cmd.Dir = cfg.ServerPath

	// 使用 execAndStream 的逻辑，但这里是后台运行，所以我们需要捕获输出并在后台推送到 WebSocket
	// 不能使用 cmd.Wait() 阻塞主线程响应
	// 我们可以启动两个 goroutine 分别读取 stdout 和 stderr
	
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fail(c, "启动失败: "+err.Error())
		return
	}

	serverProcess = cmd
	ws.Broadcast("游戏服务端正在启动...")

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
			ws.Broadcast("SERVER ERROR: " + scanner.Text())
		}
	}()

	// 异步等待进程结束
	go func() {
		cmd.Wait()
		serverMu.Lock()
		serverProcess = nil
		serverMu.Unlock()
		ws.Broadcast("游戏服务端已停止。")
	}()

	success(c, map[string]int{"pid": cmd.Process.Pid})
}

// StopServer 停止服务端
func StopServer(c *gin.Context) {
	serverMu.Lock()
	defer serverMu.Unlock()

	if serverProcess == nil || serverProcess.Process == nil {
		fail(c, "服务端未运行")
		return
	}

	ws.Broadcast("正在停止游戏服务端...")

	var err error
	if runtime.GOOS == "windows" {
		err = serverProcess.Process.Kill()
	} else {
		err = serverProcess.Process.Signal(syscall.SIGTERM)
	}

	if err != nil {
		fail(c, "停止失败: "+err.Error())
		return
	}

	success(c, nil)
}

// RestartServer 重启服务端
func RestartServer(c *gin.Context) {
	ws.Broadcast("正在重启游戏服务端...")
	serverMu.Lock()

	if serverProcess != nil && serverProcess.Process != nil {
		if runtime.GOOS == "windows" {
			serverProcess.Process.Kill()
		} else {
			serverProcess.Process.Signal(syscall.SIGTERM)
		}
		serverProcess.Wait()
		serverProcess = nil
	}

	serverMu.Unlock()

	// 启动新实例
	StartServer(c)
}
