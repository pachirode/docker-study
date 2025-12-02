package fs2

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
)

var SubsystemIns = []resource.Subsystem{
	&CpusetSubSystem{},
	&MemorySubSystem{},
	&CpuSubSystem{},
}
