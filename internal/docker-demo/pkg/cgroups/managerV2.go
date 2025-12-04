package cgroups

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/fs2"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
)

type CgroupManagerV2 struct {
	Path string
}

var _ CgroupManager = &CgroupManagerV2{}

func NewCgroupManagerV2(path string) *CgroupManagerV2 {
	return &CgroupManagerV2{
		Path: path,
	}
}

// Apply 将进程pid加入到这个cgroup中
func (c *CgroupManagerV2) Apply(pid int, res *resource.ResourceConfig) error {
	for _, subSysIns := range fs2.SubsystemIns {
		subSysIns.Apply(c.Path, pid, res)
	}
	return nil
}

// Set 设置cgroup资源限制
func (c *CgroupManagerV2) Set(res *resource.ResourceConfig) error {
	for _, subSysIns := range fs2.SubsystemIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// Destroy 释放cgroup
func (c *CgroupManagerV2) Destroy() {
	for _, subSysIns := range fs2.SubsystemIns {
		subSysIns.Remove(c.Path)
	}
}
