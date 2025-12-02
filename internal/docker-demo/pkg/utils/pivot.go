package utils

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

var (
	pivotOnce sync.Once
)

func InitBusyboxRoot() {
	pivotOnce.Do(func() {
		if ok, _ := PathExists(consts.BUSYBOX_ROOT); !ok {
			MkdirAll(consts.BUSYBOX_ROOT, consts.PERM_0777)
			runCommand("docker", "run", "-d", "--name", "my_busybox", "busybox", "top", "-b")

			runCommand("docker", "export", "-o", "busybox.tar", "my_busybox")

			runCommand("tar", "-xvf", "busybox.tar", "-C", "/root/busybox/")
		} else {
			fmt.Println("ok")
		}
	})
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\nOutput: %s", err, output)
	}
	fmt.Printf("Output of %s: %s\n", command, output)
}
