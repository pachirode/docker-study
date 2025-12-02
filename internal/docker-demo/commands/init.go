package commands

import (
	"sync"

	"github.com/pachirode/docker-demo/pkg/app"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
)

const iniCommandDesc = `Creat a container process run user's process in container`

var (
	iniOnce sync.Once
	initCmd *app.Command
)

func NewInitCommand() *app.Command {
	iniOnce.Do(func() {
		initCmd = app.NewCommand("init",
			iniCommandDesc,
			app.WithRunCommandFunc(ini()),
		)

	})
	return initCmd
}

func ini() app.RunCommandFunc {
	return func(args []string) error {

		return container.RunContainerInitProcess()
	}
}
