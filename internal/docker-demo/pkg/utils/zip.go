package utils

import (
	"os/exec"

	"github.com/pachirode/pkg/log"
)

func UnTar(tarPath, dirPath string) error {
	if _, err := exec.Command("tar", "-xvf", tarPath, "-C", dirPath).CombinedOutput(); err != nil {
		log.Errorw(err, "Error to tar file", "tarPath", tarPath)
		return err
	}

	return nil
}
