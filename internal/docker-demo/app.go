package docker_demo

import (
	"github.com/pachirode/docker-demo/pkg/app"

	"github.com/pachirode/docker-demo/internal/docker-demo/commands"
)

const commandDesc = `A simple runtime implementation.
The purpose of this project is to learn docker
`

func NewApp(basename string) *app.App {
	application := app.NewApp("Docker demo test",
		basename,
		app.WithDescription(commandDesc),
		app.WithDefaultValidArgs(),
	)

	application.AddCommand(commands.NewRunCommand())
	application.AddCommand(commands.NewInitCommand())
	application.AddCommand(commands.NewCommitCommand())
	application.AddCommand(commands.NewPSCommand())
	application.AddCommand(commands.NewLogsCommand())
	application.AddCommand(commands.NewExecCommand())
	application.AddCommand(commands.NewStopCommand())
	application.AddCommand(commands.NewRmCommand())
	application.AddCommand(commands.NewNetWorkCommand())
	application.BuildCommand()

	return application
}
