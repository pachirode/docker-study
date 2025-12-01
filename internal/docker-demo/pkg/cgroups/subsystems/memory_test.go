package subsystems

import (
	"os"
	"path"
	"testing"
)

func TestMemoryCgroup(t *testing.T) {
	memSubSystem := MemorySubSystem{}
	resConfig := &ResourceConfig{
		MemoryLimit: "1000m",
	}
	testGroup := "testmemorylimit"

	if err := memSubSystem.Set(testGroup, resConfig); err != nil {
		t.Fatalf("cgroup stats %v", err)
	}

	stat, _ := os.Stat(path.Join(FindCgroupMountpoint("memory"), testGroup))
	t.Logf("cgroup stats: %+v", stat)

	if err := memSubSystem.Apply(testGroup, os.Getpid(), resConfig); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}

	if err := memSubSystem.Apply("", os.Getpid(), resConfig); err != nil {
		t.Fatalf("cgroup Apply %v", err)
	}

	if err := memSubSystem.Remove(testGroup); err != nil {
		t.Fatalf("cgroup remove %v", err)
	}
}
