package fs

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

type CpusetSubSystem struct{}

var _ resource.Subsystem = &CpusetSubSystem{}

func (s *CpusetSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), consts.PERM_0644); err != nil {
				return fmt.Errorf("set cgroup cpuset fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int, res *resource.ResourceConfig) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
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

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}
