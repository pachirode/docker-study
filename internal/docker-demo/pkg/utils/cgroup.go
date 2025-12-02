package utils

import (
	"bufio"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/pachirode/docker-demo/pkg/errors"
	"golang.org/x/sys/unix"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

const cgroupMountPoint = "/sys/fs/cgroup"

var (
	cgroupOnce sync.Once
	isCgroup2  bool
	cgroupRoot string
)

func IsCgroup2() bool {
	cgroupOnce.Do(func() {
		var st unix.Statfs_t
		err := unix.Statfs(cgroupMountPoint, &st)
		if err != nil && os.IsNotExist(err) {
			// 对于 rootless 容器，忽略错误
			isCgroup2 = false
			return
		}
		isCgroup2 = st.Type == unix.CGROUP2_SUPER_MAGIC
	})
	return isCgroup2
}

func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text() // 104 85 0:20 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime - cgroup cgroup rw,memory
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
	if IsCgroup2() {
		cgroupRoot = cgroupMountPoint
	} else {
		cgroupRoot = FindCgroupMountPoint(subsystem)
	}

	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreate {
		return absPath, nil
	}
	// 指定自动创建时才判断是否存在
	_, err := os.Stat(absPath)
	// 只有不存在才创建
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(absPath, consts.PERM_0755)
		return absPath, err
	}
	return absPath, errors.Wrap(err, "create cgroup")
}
