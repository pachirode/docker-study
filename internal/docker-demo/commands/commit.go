package commands

import (
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs"
	"github.com/pachirode/docker-demo/pkg/app"
)

const commitCommandDesc = `commit container to image`

var (
	commitOnce sync.Once
	commitCmd  *app.Command
)

func NewCommitCommand() *app.Command {
	commitOnce.Do(func() {
		commitCmd = app.NewCommand("commit",
			commitCommandDesc,
			app.WithRunCommandFunc(commit()),
		)
	})
	return commitCmd
}

func commit() app.RunCommandFunc {
	return func(args []string) error {
		rootfs.CommitContainer()
		return nil
	}
}
