package cgroups

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

type CgroupManager interface {
	Apply(pid int, res *resource.ResourceConfig) error
	Set(res *resource.ResourceConfig) error
	Destroy() error
}

func NewCgroupManager(path string) CgroupManager {
	if utils.IsCgroup2() {
		return NewCgroupManagerV2(path)
	}

	return NewCgroupManagerV1(path)
}
