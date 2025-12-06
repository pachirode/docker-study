package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/pkg/app"
)

const networkCommandDesc = `container network commands`

var (
	networkOnce sync.Once
	networkCmd  *app.Command
	networkOpts *options.NetworkOptions
)

func NewNetWorkCommand() *app.Command {
	networkOnce.Do(func() {
		networkOpts = options.NewNetworkOptions()
		networkCmd = app.NewCommand("network",
			networkCommandDesc,
			app.WithCommandOption(networkOpts),
			app.WithRunCommandFunc(Network()),
		)
	})

	return networkCmd
}

func Network() app.RunCommandFunc {
	return func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("missing network")
		}

		op := args[0]
		switch op {
		case "create":
		case "remove":
		case "list":
		default:
			fmt.Errorf("Error op command")
		}

		return nil
	}
}
