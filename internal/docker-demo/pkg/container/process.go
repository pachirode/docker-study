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
func NewParentProcess(imageName, containID string, opts *options.RunOptions) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorw(err, "Error to create new pipe")
		return nil, nil
	}

	cmd := exec.Command(consts.BASH, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if opts.TTY {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		stdLogFile, err := GetLogFile(containID)
		if err != nil {
			log.Errorw(err, "Error to get log file")
			return nil, nil
		}
		cmd.Stdout = stdLogFile
		cmd.Stderr = stdLogFile
	}

	cmd.Env = append(os.Environ(), opts.Envs...)
	cmd.ExtraFiles = []*os.File{readPipe} // 创建的时候外带文件句柄去创建子进程
	rootfs.NewWorkSpace(containID, imageName, opts.Volume)
	cmd.Dir = rootfs.GetMerged(containID)
	return cmd, writePipe // 返回给父进程用来写入数据
}
