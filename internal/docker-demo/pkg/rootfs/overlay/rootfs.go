package overlay

import (
	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

// NewWorkSpace 创建一个 overlay 作为容器的根目录
func NewWorkSpace(opts *options.RunOptions) {
	createLower()
	createOtherDirs(consts.BUSYBOX)
	mountOverlayFS(consts.BUSYBOX)

	if opts.Volume != "" {
		mntPath := GetMerged(consts.BUSYBOX)
		hostPath, containerPath, err := volumeExtract(opts.Volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume", "volume", opts.Volume)
			return
		}
		mountVolume(mntPath, hostPath, containerPath)
	}
}

// DeleteWorkSpace 删除指定文件系统，如果容器存在
func DeleteWorkSpace(opts *options.RunOptions) {
	if opts.Volume != "" {
		mntPath := GetMerged(consts.BUSYBOX)
		_, containerPath, err := volumeExtract(opts.Volume)
		if err != nil {
			log.Errorw(err, "Error to extract volume", "volume", opts.Volume)
			return
		}
		unmountVolume(mntPath, containerPath)
	}

	unmountOverlayFS(consts.BUSYBOX)
	deleteAllDirs(consts.BUSYBOX)
}
