package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), consts.PERM_0755); err != nil {
				log.Errorw(err, "Error to create cgroup dir", "root", cgroupRoot, "path", cgroupPath)
				return "", fmt.Errorf("Error to create cgroup dir %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		log.Errorw(err, "Error to get cgroup path", "root", cgroupRoot, "path", cgroupPath)
		return "", fmt.Errorf("Error to get cgroup path %v", err)
	}
}
