package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct{}

var _ Subsystem = &MemorySubSystem{}

func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0o644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.RemoveAll(subSystemCgroupPath)
	} else {
		return err
	}
}

func (s *MemorySubSystem) Apply(cgroupPath string, pid int, res *ResourceConfig) error {
	if subSystemCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if res.MemoryLimit != "" {
			if err = os.WriteFile(path.Join(subSystemCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
				return fmt.Errorf("set cgroup proc fail %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (s *MemorySubSystem) Name() string {
	return "memory"
}
