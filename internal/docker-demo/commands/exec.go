package commands

import (
	"fmt"
	"os"
	"sync"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	_ "github.com/pachirode/docker-demo/internal/docker-demo/pkg/nsenter"
	"github.com/pachirode/docker-demo/pkg/app"
)

const execCommandDesc = `exec a command into container`

var (
	execOnce sync.Once
	execCmd  *app.Command
)

func NewExecCommand() *app.Command {
	execOnce.Do(func() {
		execCmd = app.NewCommand("exec",
			execCommandDesc,
			app.WithRunCommandFunc(exec()),
		)
	})
	return execCmd
}

func exec() app.RunCommandFunc {
	return func(args []string) error {
		if os.Getenv(consts.ENV_EXEC_PID) != "" {
			log.Infow("Pid callback", "pid", os.Getgid())
			return nil
		}

		if len(args) < 2 {
			return fmt.Errorf("Missing containerID or commands")
		}

		container.ExecContainer(args[0], args[1:])
		return nil
	}
}
