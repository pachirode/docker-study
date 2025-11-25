//go:build linux
// +build linux

package main

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Cloneflags: syscall.CLONE_NEWUTS,
		//Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC,
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		slog.Error("Error to run cmd", "err", err)
	}
}
