package fs

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/cgroups/resource"
	"os"
	"path"
	"testing"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func TestMemoryCgroup(t *testing.T) {
	memSubSystem := MemorySubSystem{}
	resConfig := &resource.ResourceConfig{
		MemoryLimit: "1000m",
	}
	testGroup := "testmemorylimit"

	if err := memSubSystem.Set(testGroup, resConfig); err != nil {
		t.Fatalf("cgroup stats %v", err)
	}

	stat, _ := os.Stat(path.Join(utils.FindCgroupMountPoint("memory"), testGroup))
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
