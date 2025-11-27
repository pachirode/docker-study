package docker_demo

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/commands"
	"github.com/pachirode/docker-demo/pkg/app"
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
	application.BuildCommand()

	return application
}
