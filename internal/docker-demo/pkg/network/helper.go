package network

import (
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

var (
	drivers = map[string]Driver{}

	ipAllocator = &IPAM{
		SubnetAllocatorPath: consts.NETWORK_DEFAULT_IPAM_ALLOCTOR,
	}
)

func init() {
	var bridgeDriver = BridgeNetworkDriver{}
	drivers[bridgeDriver.Name()] = &bridgeDriver

	utils.MkdirDirs([]string{consts.NETWORK_IPAM, consts.NETWORK_ROOT}, consts.PERM_0644)
}
