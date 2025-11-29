package utils

import (
	"os"

	"github.com/pachirode/pkg/log"
)

func MkdirAll(dirPath string, perm os.FileMode) error {
	if err := os.MkdirAll(dirPath, perm); err != nil {
		log.Errorw(err, "Error to create dir", "dirPath", dirPath, "mode", perm)
		return err
	}
	return nil
}

func MkdirDirs(dirPaths []string, perm os.FileMode) {
	for _, dirPath := range dirPaths {
		if err := os.Mkdir(dirPath, perm); err != nil {
			log.Errorw(err, "Error to create dir", "dirPath", dirPath)
		}
	}
}

func RemoveDirs(dirPaths []string) {
	for _, dirPath := range dirPaths {
		if err := os.RemoveAll(dirPath); err != nil {
			log.Errorf("Remove dir %s error %v", dirPath, err)
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
