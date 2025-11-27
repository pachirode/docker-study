package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	docker_demo "github.com/pachirode/docker-demo/internal/docker-demo"
)

func main() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	docker_demo.NewApp("docker-demo").Run()
}
