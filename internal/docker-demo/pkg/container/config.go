package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/errors"
)

const (
	RUNNING = "running"
	STOP    = "stopped"
	EXIT    = "exited"
)

type Info struct {
	Pid         string `json:"pid"`          // 容器的init进程在宿主机上的 PID
	Id          string `json:"id"`           // 容器Id
	Name        string `json:"name"`         // 容器名
	Command     string `json:"command"`      // 容器内init运行命令
	CreatedTime string `json:"created_time"` // 创建时间
	Status      string `json:"status"`       // 容器的状态
	Volume      string `json:"volume"`       // 挂载目录
}

func RecordContainerInfo(containerPID int, commandArray []string, containerName, containerID, volume string) (*Info, error) {
	if containerName == "" {
		containerName = containerID
	}
	command := strings.Join(commandArray, " ")
	containerInfo := &Info{
		Pid:         strconv.Itoa(containerPID),
		Id:          containerID,
		Name:        containerName,
		Command:     command,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
		Volume:      volume,
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
