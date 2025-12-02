package fs

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
)

var SubsystemsIns = []resource.Subsystem{
	&CpusetSubSystem{},
	&MemorySubSystem{},
	&CpuSubSystem{},
}
