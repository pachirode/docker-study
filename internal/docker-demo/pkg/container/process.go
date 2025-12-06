package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

// NewParentProcess 构建 command 用于启动第一个新的进程
func NewParentProcess(containerID string, opts *options.RunOptions, command, imageName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorw(err, "Error to create pipe")
		return nil, nil
	}
	args := []string{"init", command}
	cmd := exec.Command(consts.SELF_EXE, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if opts.TTY {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		logDir := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
		utils.MkdirAll(logDir, consts.PERM_0622)
		stdLogFilePath := path.Join(logDir, consts.CONTAINER_LOG)
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorw(err, "Error to create container log file", "path", stdLogFilePath)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(), opts.Envs...)
	rootfs.NewWorkSpace(opts, imageName)
	cmd.Dir = fmt.Sprintf(consts.LOWER_DIR_TEMP, imageName)

	return cmd, writePipe
}
