package cgroups

import (
	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/fs"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
)

type CgroupManagerV1 struct {
	// cgroup在hierarchy中的路径 相当于创建的 cgroup 目录相对于 root cgroup 目录的路径
	Path string
}

var _ CgroupManager = &CgroupManagerV1{}

func NewCgroupManagerV1(path string) *CgroupManagerV1 {
	return &CgroupManagerV1{
		Path: path,
	}
}

// Apply 将 pid 添加到 cgroup
func (c *CgroupManagerV1) Apply(pid int, res *resource.ResourceConfig) error {
	for _, subSysIns := range fs.SubsystemsIns {
		subSysIns.Apply(c.Path, pid, res)
	}
	return nil
}

// Set 设置 cgroup 资源
func (c *CgroupManagerV1) Set(res *resource.ResourceConfig) error {
	for _, subSysIns := range fs.SubsystemsIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// Destroy 删除 cgroup
func (c *CgroupManagerV1) Destroy() {
	for _, subSysIns := range fs.SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			log.Errorw(err, "Error to remove cgroup")
		}
	}
}
