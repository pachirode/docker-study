package container

import (
	"fmt"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func RemoveContainer(containerID string) {
	containerInfo, err := getContainerInfo(containerID)
	if err != nil {
		log.Errorw(err, "Error to get container info", "id", containerID)
		return
	}
	if containerInfo.Status != STOP {
		log.Errorf("Couldn't remove running container")
		return
	}
	dirDir := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	utils.RemoveDirs([]string{dirDir})
	opts := options.RunOptions{}
	opts.Volume = containerInfo.Volume
	rootfs.DeleteWorkSpace(&opts, containerInfo.Image)
}
