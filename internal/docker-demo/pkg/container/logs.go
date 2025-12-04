package container

import (
	"fmt"
	"os"
	"path"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

func GetContainerLogs(containerID string) {
	logDir := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	logPath := path.Join(logDir, consts.CONTAINER_LOG)

	context, err := os.ReadFile(logPath)
	if err != nil {
		log.Errorw(err, "Error to read container log", "path", logPath)
		return
	}
	_, err = fmt.Fprintf(os.Stdout, string(context))
	if err != nil {
		log.Errorw(err, "Error to print container log")
	}
}
