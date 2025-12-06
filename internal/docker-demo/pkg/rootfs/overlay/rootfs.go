package overlay

import (
	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
)

// NewWorkSpace 创建一个 overlay 作为容器的根目录
func NewWorkSpace(opts *options.RunOptions, imageName string) {
	createLower(imageName)
	createOtherDirs(imageName)
	mountOverlayFS(imageName)

	if opts.Volume != "" {
		mntPath := GetMerged(imageName)
		hostPath, containerPath, err := volumeExtract(opts.Volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume", "volume", opts.Volume)
			return
		}
		mountVolume(mntPath, hostPath, containerPath)
	}
}

// DeleteWorkSpace 删除指定文件系统，如果容器存在
func DeleteWorkSpace(opts *options.RunOptions, imageName string) {
	if opts.Volume != "" {
		mntPath := GetMerged(imageName)
		_, containerPath, err := volumeExtract(opts.Volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume", "volume", opts.Volume)
			return
		}
		unmountVolume(mntPath, containerPath)
	}

	unmountOverlayFS(imageName)
	deleteAllDirs(imageName)
}
