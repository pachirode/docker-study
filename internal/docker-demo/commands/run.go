package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/pkg/app"
)

const commandDesc = `Create a container with namespace and cgroups limit
docker-demo run -it [command]
docker-demo run -d -name [containerName] [imageName] [command]
`

var (
	once    sync.Once
	runCmd  *app.Command
	runOpts *options.RunOptions
)

func NewRunCommand() *app.Command {

	once.Do(func() {
		runOpts = options.NewRunOptions()
		runCmd = app.NewCommand("run",
			commandDesc,
			app.WithCommandOption(runOpts),
			app.WithRunCommandFunc(Run()),
		)
	})

	return runCmd
}

func Run() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing container command")
		}

		if runOpts.TTY && runOpts.Detach {
			return fmt.Errorf("tty and detach flag can not both provided")
		}
		if !runOpts.Detach {
			runOpts.TTY = true
		}

		var cmdArray []string
		for _, arg := range args {
			cmdArray = append(cmdArray, arg)
		}
		//imagName := cmdArray[0]

		log.Infow("Starting Create tty")
		//run(imagName, runOpts)
		return nil
	}
}

//func run(imagName string, opts *options.RunOptions) {
//	containerID := container.GenerateContainerID()
//
//	// Start
//	parent, writePipe := container.NewParentProcess(imagName, containerID, opts)
//	if parent == nil {
//		log.Errorf("Error to create parent process")
//		return
//	}
//	if err := parent.Start(); err != nil {
//		log.Errorw(err, "Error to run parent.Start")
//		return
//	}
//}
