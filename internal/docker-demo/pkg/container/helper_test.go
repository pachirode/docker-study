package container

import "testing"

func TestGenerateContainerID(t *testing.T) {
	str, err := GenerateContainerID()
	if err != nil {
		t.Errorf("Error to generate containerID, err: %v", err)
	}
	if len(str) != 12 {
		t.Errorf("generate contianerID length err, want 12, got %d", len(str))
	}
}
