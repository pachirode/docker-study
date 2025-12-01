package container

import (
	"os"
	"syscall"

	"github.com/pachirode/pkg/log"
)

// RunContainerInitProcess 启动容器的 init 进程
func RunContainerInitProcess(command string, arg []string) error {
	// 必须声明新的 mount namespace 独立
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Errorw(err, "Error to run container init process")
	}
	return nil
}
