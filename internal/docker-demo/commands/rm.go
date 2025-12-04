package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	"github.com/pachirode/docker-demo/pkg/app"
)

const rmCommandDesc = `rm a container`

var (
	rmOnce sync.Once
	rmCmd  *app.Command
)

func NewRmCommand() *app.Command {
	rmOnce.Do(func() {
		rmCmd = app.NewCommand("rm",
			rmCommandDesc,
			app.WithRunCommandFunc(rm()),
		)

	})
	return rmCmd
}

func rm() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing containerID")
		}
		container.RemoveContainer(args[0])
		return nil
	}
}
