package container

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

func ListContainers() {
	files, err := os.ReadDir(consts.INFO_LOCATION)
	if err != nil {
		log.Errorw(err, "Error to list files", "path", consts.INFO_LOCATION)
		return
	}
	containers := make([]*Info, 0)
	for _, file := range files {
		if file.IsDir() {
			tmpContainer, err := getContainerInfo(file.Name())
			if err != nil {
				log.Errorw(err, "Error to get container info", "path", file)
				continue
			}
			containers = append(containers, tmpContainer)
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 14, 1, 3, ' ', 0)
	_, err = fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	if err != nil {
		log.Errorw(err, "Error to set print template")
		return
	}
	for _, item := range containers {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreatedTime,
		)
		if err != nil {
			log.Errorw(err, "Error to format config info", "id", item.Id)
		}
	}
	if err = w.Flush(); err != nil {
		log.Errorw(err, "Error to flush container infos")
	}
}
