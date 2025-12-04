package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

// RunContainerInitProcess 启动容器的 init 进程
func RunContainerInitProcess() error {
	cmdArray := utils.ReadPipeCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Error to read pip command, command is nil")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorw(err, "Error to exec loop path", "cmd", cmdArray[0])
	}

	if err = syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorw(err, "Error to run container init process")
	}
	return nil
}
