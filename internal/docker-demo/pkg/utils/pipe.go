package utils

import (
	"io"
	"os"
	"strings"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

// WritePipeCommand 通过 WritePipe 将指令发送
func WritePipeCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infow("WritePipe command msg", "cmd", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}

// ReadPipeCommand 通过 ReadPipe 读取指令
func ReadPipeCommand() []string {
	pipe := os.NewFile(uintptr(consts.FDINDEX), "pipe")
	defer pipe.Close()
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Errorw(err, "Error to read pipe")
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
