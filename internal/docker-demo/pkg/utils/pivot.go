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

func InitImageRoot(imageName string) {
	pivotOnce.Do(func() {
		imageRoot := fmt.Sprintf(consts.LOWER_DIR_TEMP, strings.Replace(imageName, ":", "_", -1))
		if ok, _ := PathExists(imageRoot); !ok {
			MkdirAll(imageRoot, consts.PERM_0777)
			runCommand("docker", "run", "-d", "--name", "my_busybox", imageName, "top", "-b")

			runCommand("docker", "export", "-o", "image.tar", "my_busybox")

			runCommand("docker", "stop", "my_busybox")

			runCommand("docker", "rm", "my_busybox")

			runCommand("tar", "-xvf", "image.tar", "-C", imageRoot)

			os.Remove("image.tar")
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
