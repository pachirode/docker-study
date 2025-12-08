package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/config"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/errors"
)

const (
	RUNNING = "running"
	STOP    = "stopped"
	EXIT    = "exited"
)

func RecordContainerInfo(containerPID int, commandArray []string, containerName, containerID, volume, imageName, networkName, ip string, portMapping []string) (*config.Info, error) {
	if containerName == "" {
		containerName = containerID
	}
	command := strings.Join(commandArray, " ")
	containerInfo := &config.Info{
		Pid:         strconv.Itoa(containerPID),
		Id:          containerID,
		Name:        containerName,
		Command:     command,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
		Volume:      volume,
		Image:       imageName,
		NetworkName: networkName,
		PortMapping: portMapping,
		IP:          ip,
	}
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		return containerInfo, errors.WithMessage(err, "Error to marshal container info")
	}
	jsonStr := string(jsonBytes)
	dirPath := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	utils.MkdirAll(dirPath, consts.PERM_0644)
	fileName := path.Join(dirPath, consts.CONFIG_JSON)
	file, err := os.Create(fileName)
	if err != nil {
		return containerInfo, errors.WithMessage(err, "Error to create file")
	}
	defer file.Close()
	if _, err = file.WriteString(jsonStr); err != nil {
		return containerInfo, errors.WithMessage(err, "Error to container info")
	}
	return containerInfo, nil
}

func DeleteContainerInfo(containerID string) {
	dirPath := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	utils.RemoveDirs([]string{dirPath})
}
