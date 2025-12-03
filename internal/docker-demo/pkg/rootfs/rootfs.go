package rootfs

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs/overlay"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func NewWorkSpace(opts *options.RunOptions) {
	if utils.SupportOverlay() {
		overlay.NewWorkSpace(opts)
	}
}

func DeleteWorkSpace(opts *options.RunOptions) {
	if utils.SupportOverlay() {
		overlay.DeleteWorkSpace(opts)
	}
}

func CommitContainer() {
	if utils.SupportOverlay() {
		overlay.CommitContainer()
	}
}
