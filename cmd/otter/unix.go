//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// https://stackoverflow.com/questions/22470193/why-wont-go-kill-a-child-process-correctly
func stop(cmd *exec.Cmd) {
	pgid := -cmd.Process.Pid
	_ = syscall.Kill(pgid, syscall.SIGTERM)
}

func setpgid(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
