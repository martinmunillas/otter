package main

import (
	"log/slog"
	"os"
	"os/exec"
)

func createDefaultCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func fatal(logger *slog.Logger, err error) {
	logger.Error(err.Error())
	os.Exit(1)
}
