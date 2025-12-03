package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs"
)

// NewParentProcess 构建 command 用于启动第一个新的进程
func NewParentProcess(opts *options.RunOptions, command string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorw(err, "Error to create pipe")
		return nil, nil
	}
	args := []string{"init", command}
	cmd := exec.Command(consts.BASH, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if opts.TTY {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	rootfs.NewWorkSpace(opts)
	cmd.Dir = consts.BUSYBOX_ROOT_MNT

	return cmd, writePipe
}
