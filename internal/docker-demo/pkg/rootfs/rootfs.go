package rootfs

import "github.com/pachirode/pkg/log"

// NewWorkSpace 创建一个 overlay 作为容器的根目录
func NewWorkSpace(containerID, imageName, volume string) {
	createLower(containerID, imageName)
	createOtherDirs(containerID)
	mountOverlayFS(containerID)

	if volume != "" {
		mntPath := GetMerged(containerID)
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume")
			return
		}
		mountVolume(mntPath, hostPath, containerPath)
	}
}

// DeleteWorkSpace 删除指定文件系统，如果容器存在
func DeleteWorkSpace(containerID, volume string) {
	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume")
			return
		}
		mntPath := GetMerged(containerID)
		unmountVolume(mntPath, containerPath)
	}

	unmountOverlayFS(containerID)
	deleteAllDirs(containerID)
}
