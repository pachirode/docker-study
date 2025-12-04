package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	"github.com/pachirode/docker-demo/pkg/app"
)

const logsCommandDesc = `Creat a container process run user's process in container`

var (
	logsOnce sync.Once
	logsCmd  *app.Command
)

func NewLogsCommand() *app.Command {
	logsOnce.Do(func() {
		logsCmd = app.NewCommand("logs",
			logsCommandDesc,
			app.WithRunCommandFunc(logs()),
		)

	})
	return logsCmd
}

func logs() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing containerID")
		}
		container.GetContainerLogs(args[0])
		return nil
	}
}
