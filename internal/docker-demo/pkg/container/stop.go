package container

import (
	"encoding/json"
	"fmt"
	"github.com/pachirode/docker-demo/pkg/errors"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

func StopContainer(containerID string) {
	pid, err := getContainerPidByID(containerID)
	if err != nil {
		log.Errorw(err, "Error to get container pid", "id", containerID)
		return
	}
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Errorw(err, "Error to convert pid from string to int", "pid", pid)
		return
	}
	p, err := os.FindProcess(pidInt)
	if err != nil {
		log.Errorw(err, "Error to get process by pid", "pid", pid)
		return
	}
	if err = p.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		log.Errorw(err, "Error to stop container", "id", containerID, "pid", pid)
		return
	}
	info, err := getContainerInfo(containerID)
	if err != nil {
		log.Errorw(err, "Error to get container info", "id", containerID)
		return
	}
	info.Status = STOP
	info.Pid = "-1"
	newContext, err := json.Marshal(info)
	if err != nil {
		log.Errorw(err, "Error to marshal container info", "data", info)
		return
	}
	containerDir := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	configPath := path.Join(containerDir, consts.CONFIG_JSON)
	if err = os.WriteFile(configPath, newContext, consts.PERM_0622); err != nil {
		log.Errorw(err, "Error to update container info", "id", containerID)
	}
}
