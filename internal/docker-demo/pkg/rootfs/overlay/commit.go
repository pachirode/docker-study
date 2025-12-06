package overlay

import (
	"fmt"
	"os/exec"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func CommitContainer(imageName string) {
	imageTar := GetImage(imageName)
	exists, err := utils.PathExists(imageTar)
	if err != nil {
		log.Errorw(err, "Error to check image", "imagePath", imageTar)
		return
	}

	if exists {
		log.Infow("Image exists")
		return
	}

	log.Infow("Starting commit container", "imagePath", imageTar)
	imageRoot := fmt.Sprintf(consts.LOWER_DIR_TEMP, imageName)
	if _, err = exec.Command("tar", "-czf", imageTar, "-C", imageRoot, ".").CombinedOutput(); err != nil {
		log.Errorw(err, "Error to tar container", "mnt", imageRoot)
	}
}
