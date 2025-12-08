package container

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/config"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/pkg/errors"
)

func GenerateContainerID() string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GetInfoLocation(containerID string) (string, error) {
	dirPath := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	if err := os.MkdirAll(dirPath, consts.PERM_0622); err != nil {
		return "", err
	}

	return dirPath, nil
}

func GetLogFile(containerID string) (*os.File, error) {
	GetInfoLocation(containerID)
	stdLogFilePath := fmt.Sprintf(consts.LOG_FILE_TEMP, containerID, containerID)
	stdLogFile, err := os.Create(stdLogFilePath)
	if err != nil {
		return nil, err
	}
	return stdLogFile, nil
}

// Init 挂载
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorw(err, "Error to get location")
		return
	}
	log.Infow("Get current location successfully", "dirPath", pwd)

	// ""没有特定的源文件系统，"/" 目标挂载点的根文件系统，私有挂载，递归有效
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	err = pivotRoot(pwd)
	if err != nil {
		log.Errorw(err, "Error to pivotRoot")
		return
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// 挂载 proc 文件系统，通常用户访问内核和进程信息
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

// pivotRoot 实现系统调用功能，将当前进程的根文件系统切换到新的根文件系统
// 实现隔离和沙盒
func pivotRoot(root string) error {
	// 确保新旧根文件系统挂载到自己，将同一个文件系统挂载到不同的挂载点
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return errors.Wrap(err, "Error to mount rootfs to itself")
	}
	// pivot_root 要求新根和旧根不在同一个目录，需要新创建目录
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, consts.PERM_0777); err != nil {
		return err
	}
	// 将根文件切换到新的 root
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return errors.WithMessagef(err, "Error to pivotRoot")
	}

	if err := syscall.Chdir("/"); err != nil {
		return errors.WithMessage(err, "Error to chdir to /")
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// 卸载旧的
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return errors.WithMessage(err, "Error to unmount pivot_root")
	}

	return os.Remove(pivotDir)
}

func getContainerInfo(containerID string) (*config.Info, error) {
	configDir := fmt.Sprintf(consts.INFO_LOCATION_TEMP, containerID)
	configFile := path.Join(configDir, consts.CONFIG_JSON)
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Errorw(err, "Error to read config file", "path", configFile)
		return nil, err
	}
	info := new(config.Info)
	if err = json.Unmarshal(content, info); err != nil {
		log.Errorw(err, "Error to unmarshal config file", "content", content)
		return nil, err
	}

	return info, nil
}

func getContainerPidByID(containerID string) (string, error) {
	info, err := getContainerInfo(containerID)
	if err != nil {
		return "", err
	}
	return info.Pid, nil
}
