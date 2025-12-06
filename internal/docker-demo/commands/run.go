package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/container"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/rootfs"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/app"
)

const runCommandDesc = `Create a container with namespace and cgroups limit
docker-demo run -it [command]
docker-demo run -d -name [containerName] [imageName] [command]
`

var (
	runOnce sync.Once
	runCmd  *app.Command
	runOpts *options.RunOptions
)

func NewRunCommand() *app.Command {
	runOnce.Do(func() {
		runOpts = options.NewRunOptions()
		runCmd = app.NewCommand("run",
			runCommandDesc,
			app.WithCommandOption(runOpts),
			app.WithRunCommandFunc(Run()),
		)
	})

	return runCmd
}

func Run() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("missing container command")
		}

		if runOpts.TTY && runOpts.Detach {
			return fmt.Errorf("tty and detach flag can not both provided")
		}
		if !runOpts.Detach {
			runOpts.TTY = true
		}
		imageName := args[0]
		var cmdArray []string
		for _, arg := range args[1:] {
			cmdArray = append(cmdArray, arg)
		}

		log.Infow("Starting Create tty")
		run(runOpts, cmdArray, imageName)
		return nil
	}
}

func run(opts *options.RunOptions, cmdArray []string, imageName string) {
	containerID := container.GenerateContainerID()

	parent, writePipe := container.NewParentProcess(containerID, opts, cmdArray[0], imageName)
	if parent == nil {
		log.Errorf("Error to create parent process")
		return
	}
	if err := parent.Start(); err != nil {
		log.Errorw(err, "Error to run parent.Start")
		return
	}

	cgroupManager := cgroups.NewCgroupManager("docker_demo")
	res := &resource.ResourceConfig{
		CpuSet:      opts.CPUSet,
		CpuShare:    opts.CPUShare,
		MemoryLimit: opts.MEM,
	}
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(res)
	_ = cgroupManager.Apply(parent.Process.Pid, res)

	_, err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, opts.Name, containerID, opts.Volume, imageName)
	if err != nil {
		log.Errorw(err, "Error to record container info")
		return
	}

	utils.WritePipeCommand(cmdArray, writePipe)
	if opts.TTY {
		_ = parent.Wait()
		rootfs.DeleteWorkSpace(opts, imageName)
		container.DeleteContainerInfo(containerID)
	}
}
