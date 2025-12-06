package network

import (
	"io/fs"
	"path"
	"path/filepath"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
)

var (
	drivers  = map[string]Driver{}
	networks = map[string]*Network{}

	ipAllocator = &IPAM{
		SubnetAllocatorPath: consts.NETWORK_DEFAULT_IPAM_ALLOCTOR,
	}
)

func loadNetWork() (map[string]*Network, error) {
	err := filepath.Walk(consts.NETWORK_ROOT, func(nwPath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, netName := path.Split(nwPath)
		nw := &Network{
			Name: netName,
		}
		if err = nw.load(nwPath); err != nil {
			log.Errorw(err, "Error to load network file")
		}
		networks[netName] = nw
		return nil
	})

	return networks, err
}
