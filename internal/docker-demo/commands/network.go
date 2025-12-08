package commands

import (
	"fmt"
	"sync"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/network"
	"github.com/pachirode/docker-demo/pkg/app"
	"github.com/pachirode/docker-demo/pkg/errors"
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
		if len(args) < 1 {
			return fmt.Errorf("missing op arg")
		}

		op := args[0]
		switch op {
		case "create":
			if len(args) < 2 {
				return fmt.Errorf("missing network name")
			}
			if networkOpts.Driver == "" || networkOpts.Subnet == "" {
				log.Errorf("driver or subnet is emnpty")
				return nil
			}

			err := network.CreateNetwork(networkOpts, args[1])
			if err != nil {
				return err
			}
			return nil
		case "remove":
			if len(args) < 2 {
				return fmt.Errorf("missing network name")
			}
			err := network.DeleteNetwork(args[1])
			if err != nil {
				return errors.WithMessage(err, "Error to remove network")
			}
			return nil
		case "list":
			network.ListNetwork()
		default:
			return errors.New("Error op command, must create, list or remove")
		}

		return nil
	}
}
