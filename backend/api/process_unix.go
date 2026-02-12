//go:build !windows

package api

import "os/exec"

// hideWindow Unix 系统无需隐藏窗口
func hideWindow(cmd *exec.Cmd) {}
