package rootfs

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs/overlay"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func NewWorkSpace(opts *options.RunOptions, imageName string) {
	if utils.SupportOverlay() {
		overlay.NewWorkSpace(opts, imageName)
	}
}

func DeleteWorkSpace(opts *options.RunOptions, imageName string) {
	if utils.SupportOverlay() {
		overlay.DeleteWorkSpace(opts, imageName)
	}
}

func CommitContainer(imageName string) {
	if utils.SupportOverlay() {
		overlay.CommitContainer(imageName)
	}
}
