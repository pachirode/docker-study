package rootfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func GetRoot(containerID string) string { return consts.ROOT_PATH + containerID }

func GetImage(imageName string) string { return fmt.Sprintf("%s%s.tar", consts.IMAGE_PATH, imageName) }

func GetLower(containerID string) string {
	return fmt.Sprintf(consts.LOWER_DIR_TEMP, containerID)
}

func GetUpper(containerID string) string {
	return fmt.Sprintf(consts.UPPER_DIR_TEMP, containerID)
}

func GetWorker(containerID string) string {
	return fmt.Sprintf(consts.WORK_DIR_TEMP, containerID)
}

func GetMerged(containerID string) string { return fmt.Sprintf(consts.MERGED_DIR_TEMP, containerID) }

func GetOverlayFSDirs(lower, upper, worker string) string {
	return fmt.Sprintf(consts.OVERLAY_PARAMETER_TEMP, lower, upper, worker)
}

// createLower 根据 containerID, imageName 准备 lower 层目录
func createLower(containerID, imageName string) {
	lowerPath := GetLower(containerID)
	imagePath := GetImage(imageName)

	log.Infow("Start create lower", "dirPath", lowerPath, "imageName", imageName)

	exist, err := utils.PathExists(lowerPath)
	if err != nil {
		log.Errorw(err, "Error to judge whether dir exists", "dirPath", lowerPath)
	}

	if !exist {
		if err = utils.MkdirAll(lowerPath, consts.PERM_0777); err != nil {
			log.Errorw(err, "Error to mkdir", "dirPath", lowerPath)
			if err = utils.UnTar(imagePath, lowerPath); err != nil {
				log.Errorw(err, "Error to unTar", "image", imageName)
			}
		}
	}
}

// createOtherDirs 创建 upper, merged, worker
func createOtherDirs(containerID string) {
	dirs := []string{
		GetMerged(containerID),
		GetUpper(containerID),
		GetWorker(containerID),
	}

	utils.MkdirDirs(dirs, consts.PERM_0777)
}

// deleteAllDirs 删除所有涉及到达文件夹
func deleteAllDirs(containerID string) {
	dirs := []string{
		GetMerged(containerID),
		GetUpper(containerID),
		GetWorker(containerID),
		GetLower(containerID),
		GetRoot(containerID), // root 目录也要删除
	}

	utils.RemoveDirs(dirs)
}

// mountOverlayFS 挂载 overlay 文件系统
func mountOverlayFS(containerID string) {
	lowerPath := GetLower(containerID)
	upperPath := GetUpper(containerID)
	workerPath := GetWorker(containerID)
	mergedPath := GetMerged(containerID)

	overlayFSDirs := GetOverlayFSDirs(lowerPath, upperPath, workerPath)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", overlayFSDirs, mergedPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debugw("Start mount overlayFS")
	if err := cmd.Run(); err != nil {
		log.Errorw(err, "Error to mount overlayFS")
	}
}

// unmountOverlayFS 取消挂载目录
func unmountOverlayFS(containID string) {
	mergedPath := GetMerged(containID)
	cmd := exec.Command("unmount", mergedPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Debugw("Start unmount overlayFS", "dirPath", mergedPath)
	if err := cmd.Run(); err != nil {
		log.Errorw(err, "Error to unmount", "dirPath", mergedPath)
	}
}

// mountVolume 挂载 volume
func mountVolume(mntPath, hostPath, containerPath string) {
	err := utils.MkdirAll(hostPath, consts.PERM_0777)
	if err != nil {
		log.Errorw(err, "Error to create hostPath", "dirPath", hostPath)
		return
	}
	containerPathInHost := path.Join(mntPath, containerPath)
	if err = utils.MkdirAll(containerPathInHost, consts.PERM_0777); err != nil {
		log.Errorw(err, "Error to create containerPathInHost", "dirPath", containerPathInHost)
		return
	}
	cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Errorw(err, "Error to mount volume", "hostPath", hostPath, "mntPath", containerPathInHost)
	}
}

// unmountVolume 取消 volume 挂载
func unmountVolume(mntPath, containerPath string) {
	containerPathInHost := path.Join(mntPath, containerPath)
	cmd := exec.Command("umount", containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorw(err, "Error to umount volume")
	}
}

// volumeExtract 通过冒号分割解析 volume 目录
func volumeExtract(volume string) (sourcePath, destinationPath string, err error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Invalid volume [%s], must split by `:`", volume)
	}

	sourcePath, destinationPath = parts[0], parts[1]
	if sourcePath == "" || destinationPath == "" {
		return "", "", fmt.Errorf("Invalid volume [%s], path can't be empty", volume)
	}

	return sourcePath, destinationPath, nil
}
