//go:build windows

package api

import (
	"os/exec"
	"syscall"
)

// hideWindow 设置 Windows 隐藏控制台窗口
func hideWindow(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.HideWindow = true
}
