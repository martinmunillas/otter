//go:build windows

package main

import (
	"os/exec"
	"strconv"
)

// https://stackoverflow.com/a/44551450
func stop(cmd *exec.Cmd) {
	killCmd := exec.Command("taskkill.exe", "/t", "/f", "/pid", strconv.Itoa(cmd.Process.Pid))
	_ = killCmd.Run()
}

func setpgid(cmd *exec.Cmd) {}
