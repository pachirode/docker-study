package container

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/errors"
	"github.com/pachirode/pkg/log"
	"os"
	"os/exec"
	"syscall"
)

// RunContainerInitProcess 启动容器的 init 进程
func RunContainerInitProcess() error {
	cmdArray := utils.ReadPipeCommand()
	if len(cmdArray) == 0 {
		return errors.New("Run container get user command error")
	}

	setUpMount()

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorw(err, "Error to exec loop path", "command", cmdArray[0])
		return err
	}
	log.Infow("Find path", "path", path)
	if err = syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorw(err, "Error to run container init process")
	}
	return nil
}
