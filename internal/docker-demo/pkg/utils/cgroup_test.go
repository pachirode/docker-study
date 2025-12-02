package utils

import "testing"

func TestGetCgroupPath(t *testing.T) {
	path, err := GetCgroupPath("", "docker_demo", true)
	t.Logf("Cgroup path %v, error %v", path, err)
}
