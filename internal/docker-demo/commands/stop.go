package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	"github.com/pachirode/docker-demo/pkg/app"
)

const stopCommandDesc = `stop a container`

var (
	stopOnce sync.Once
	stopCmd  *app.Command
)

func NewStopCommand() *app.Command {
	stopOnce.Do(func() {
		stopCmd = app.NewCommand("stop",
			stopCommandDesc,
			app.WithRunCommandFunc(stop()),
		)

	})
	return stopCmd
}

func stop() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing containerID")
		}
		container.RemoveContainer(args[0])
		return nil
	}
}
