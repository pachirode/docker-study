package resource

type Subsystem interface {
	Name() string
	Set(path string, res *ResourceConfig) error
	Apply(path string, pid int, res *ResourceConfig) error
	Remove(path string) error
}

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    int
	CpuSet      string
}
