package fs2

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

type CpusetSubSystem struct {
}

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}

func (s *CpusetSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}
	subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(subCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), consts.PERM_0644); err != nil {
		return fmt.Errorf("set cgroup cpuset fail %v", err)
	}
	return nil
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int, res *resource.ResourceConfig) error {
	if res.CpuShare != 0 {
		if subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, true); err == nil {
			if err = os.WriteFile(path.Join(subCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), consts.PERM_0644); err != nil {
				log.Errorw(err, "Error to set cpu cgroup")
			}
		}
	}
	return nil
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
