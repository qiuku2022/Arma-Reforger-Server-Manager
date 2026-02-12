package api

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

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
	} else {
		executable = filepath.Join(cfg.ServerPath, "ArmaReforgerServer")
	}

	configPath := filepath.Join(cfg.ServerPath, "config.json")
	profilePath := filepath.Join(cfg.ServerPath, "profile")

	cmd = exec.Command(executable, "-config", configPath, "-profile", profilePath)
	cmd.Dir = cfg.ServerPath

	// Windows 隐藏控制台窗口
	hideWindow(cmd)

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
		// Windows: 先尝试优雅终止 (CTRL+C)，再强制终止
		err = gracefulKillWindows(serverProcess.Process.Pid)
	} else {
		// Linux: SIGTERM 优雅终止
		err = serverProcess.Process.Signal(syscall.SIGTERM)
		if err == nil {
			// 等待 3 秒后检查是否终止
			done := make(chan error, 1)
			go func() {
				done <- serverProcess.Wait()
			}()
			select {
			case <-done:
				// 已正常终止
			case <-time.After(3 * time.Second):
				// 超时，强制终止
				serverProcess.Process.Kill()
			}
		}
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
			gracefulKillWindows(serverProcess.Process.Pid)
		} else {
			serverProcess.Process.Signal(syscall.SIGTERM)
			// 等待最多 3 秒
			done := make(chan error, 1)
			go func() {
				done <- serverProcess.Wait()
			}()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				serverProcess.Process.Kill()
			}
		}
		serverProcess = nil
	}
	serverMu.Unlock()

	// 等待一小段时间确保进程完全结束
	time.Sleep(500 * time.Millisecond)

	// 启动新实例
	StartServer(c)
}

// gracefulKillWindows Windows 优雅终止进程
func gracefulKillWindows(pid int) error {
	// 尝试发送 CTRL+C 信号（优雅终止）
	// 但 Go 的 syscall 不支持直接发送 CTRL+C 到子进程
	// 使用 taskkill 的 /T 参数终止进程树
	cmd := exec.Command("taskkill", "/T", "/PID", string(rune(pid)))
	err := cmd.Run()
	if err != nil {
		// 如果 taskkill 失败，强制终止
		return exec.Command("taskkill", "/T", "/F", "/PID", string(rune(pid))).Run()
	}

	// 等待进程终止
	time.Sleep(1 * time.Second)

	// 检查进程是否还在运行
	checkCmd := exec.Command("tasklist", "/FI", "PID eq "+string(rune(pid)))
	out, _ := checkCmd.Output()
	if string(out) != "" {
		// 进程还在，强制终止
		return exec.Command("taskkill", "/T", "/F", "/PID", string(rune(pid))).Run()
	}

	return nil
}
