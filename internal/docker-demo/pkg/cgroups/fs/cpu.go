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

const (
	PeriodDefault = 100000
	Percent       = 100
)

type CpuSubSystem struct{}

var _ resource.Subsystem = &CpuSubSystem{}

func (s *CpuSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuShare != 0 {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "cpu.cfs_period_us"), []byte(strconv.Itoa(PeriodDefault)), consts.PERM_0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "cpu.cfs_quota_us"), []byte(strconv.Itoa(PeriodDefault/Percent*res.CpuShare)), consts.PERM_0644); err != nil {
				return fmt.Errorf("set cgroup cpu share fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int, res *resource.ResourceConfig) error {
	if subSystemCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if res.CpuShare != 0 {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0o644); err != nil {
				return fmt.Errorf("set cgroup proc fail %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}
