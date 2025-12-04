package commands

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	"github.com/pachirode/docker-demo/pkg/app"
	"sync"
)

const psCommandDesc = `list all the containers`

var (
	psOnce sync.Once
	psCmd  *app.Command
)

func NewPSCommand() *app.Command {
	psOnce.Do(func() {
		psCmd = app.NewCommand("ps",
			psCommandDesc,
			app.WithRunCommandFunc(ps()),
		)

	})
	return psCmd
}

func ps() app.RunCommandFunc {
	return func(args []string) error {
		container.ListContainers()
		return nil
	}
}
