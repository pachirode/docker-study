package container

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

func ExecContainer(containerID string, cmdArray []string) {
	pid, err := getContainerPidByID(containerID)
	if err != nil {
		log.Errorw(err, "Error to get container pid", "id", containerID)
		return
	}

	cmdStr := strings.Join(cmdArray, " ")
	log.Infow("Starting exec container", "pid", pid, "cmd", cmdStr)

	cmd := exec.Command(consts.SELF_EXE, "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(consts.ENV_EXEC_PID, pid)
	os.Setenv(consts.ENV_EXEC_CMD, cmdStr)

	if err = cmd.Run(); err != nil {
		log.Errorw(err, "Error to exec container")
	}
}
