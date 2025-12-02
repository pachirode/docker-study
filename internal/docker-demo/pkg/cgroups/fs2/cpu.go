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

const (
	PeriodDefault = 100000
	Percent       = 100
)

type CpuSubSystem struct {
}

var _ resource.Subsystem = &CpuSubSystem{}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}

func (s *CpuSubSystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuShare == 0 {
		return nil
	}
	subCgroupPath, err := utils.GetCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	// cpu.cfs_period_us & cpu.cfs_quota_us 控制的是CPU使用时间，单位是微秒，比如每1秒钟，这个进程只能使用200ms，相当于只能用20%的CPU
	// v2 中直接将 cpu.cfs_period_us & cpu.cfs_quota_us 统一记录到 cpu.max 中，比如 5000 10000 这样就是限制使用 50% cpu
	if res.CpuShare != 0 {
		// cpu.cfs_quota_us 则根据用户传递的参数来控制，比如参数为20，就是限制为20%CPU，所以把cpu.cfs_quota_us设置为cpu.cfs_period_us的20%就行
		if err = os.WriteFile(path.Join(subCgroupPath, "cpu.max"), []byte(fmt.Sprintf("%s %s", strconv.Itoa(PeriodDefault/Percent*res.CpuShare), PeriodDefault)), consts.PERM_0644); err != nil {
			return fmt.Errorf("set cgroup cpu share fail %v", err)
		}
	}
	return nil
}

func (s *CpuSubSystem) Apply(cgroupPath string, pid int, res *resource.ResourceConfig) error {
	if res.CpuShare != 0 {
		if subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, true); err == nil {
			if err = os.WriteFile(path.Join(subCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), consts.PERM_0644); err != nil {
				log.Errorw(err, "Error to set cpu cgroup")
			}
		}
	}
	return nil
}

func (s *CpuSubSystem) Remove(cgroupPath string) error {
	subCgroupPath, err := utils.GetCgroupPath("", cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subCgroupPath)
}
