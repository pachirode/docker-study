package fs2

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

type MemorySubSystem struct {
}

// Name 返回cgroup名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置cgroupPath对应的cgroup的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}
	subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, true)
	if err != nil {
		return err
	}
	// 设置这个cgroup的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes 文件中。
	if err = os.WriteFile(path.Join(subCgroupPath, "memory.max"), []byte(res.MemoryLimit), consts.PERM_0644); err != nil {
		return fmt.Errorf("set cgroup memory fail %v", err)
	}
	return nil
}

// Apply 将pid加入到cgroupPath对应的cgroup中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int, res *resource.ResourceConfig) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath("", cgroupPath, false); err == nil {
		if res.CpuSet != "" {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
				return fmt.Errorf("set cgroup proc fail %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

// Remove 删除cgroupPath对应的cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
