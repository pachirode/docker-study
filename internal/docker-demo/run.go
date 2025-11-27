package docker_demo

import "github.com/pachirode/docker-demo/pkg/app"

func Run() app.RunFunc {
	return func(basename string) error {
		println("Test")
		return nil
	}
}
