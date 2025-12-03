package overlay

import (
	"os/exec"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

func CommitContainer() {
	imageTar := GetImage(consts.BUSYBOX)
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
	if _, err = exec.Command("tar", "-czf", imageTar, "-C", consts.BUSYBOX_ROOT_MNT, ".").CombinedOutput(); err != nil {
		log.Errorw(err, "Error to tar container", "mnt", consts.BUSYBOX_ROOT_MNT)
	}
}
