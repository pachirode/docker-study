package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

var (
	pivotOnce      sync.Once
	overlayOnce    sync.Once
	supportOverlay bool
)

func InitBusyboxRoot() {
	pivotOnce.Do(func() {
		if ok, _ := PathExists("/var/lib/docker-demo/overlay2/busybox/lower"); !ok {
			MkdirAll("/var/lib/docker-demo/overlay2/busybox/lower", consts.PERM_0777)
			runCommand("docker", "run", "-d", "--name", "my_busybox", "busybox", "top", "-b")

			runCommand("docker", "export", "-o", "busybox.tar", "my_busybox")

			runCommand("tar", "-xvf", "busybox.tar", "-C", "/var/lib/docker-demo/overlay2/busybox/lower")
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

func SupportOverlay() bool {
	overlayOnce.Do(func() {
		file, err := os.Open("/proc/filesystems")
		if err != nil {
			log.Errorw(err, "Error to open /proc/filesystem")
			supportOverlay = false
			return
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "overlay") {
				supportOverlay = true
			}
		}
	})

	return supportOverlay
}
