package network

import (
	"encoding/json"
	"net"
	"os"
	"path"

	"github.com/pachirode/pkg/log"
	"github.com/vishvananda/netlink"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
)

type Network struct {
	Name    string
	IpRange *net.IPNet
	Driver  string
}

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

func (nw *Network) dump(dumpPath string) error {
	utils.MkdirAll(dumpPath, consts.PERM_0644)

	nwPath := path.Join(dumpPath, nw.Name)
	// 保证有一个空白的文件存在
	nwFile, err := os.OpenFile(nwPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, consts.PERM_0644)
	if err != nil {
		log.Errorw(err, "Error to open file", "path", nwPath)
		return err
	}
	defer nwFile.Close()

	nwJson, err := json.Marshal(nw)
	if err != nil {
		log.Errorw(err, "Error to marshal json", "data", nw)
		return err
	}

	_, err = nwFile.Write(nwJson)
	if err != nil {
		log.Errorw(err, "Error to write data")
		return err
	}

	return nil
}

func (nw *Network) remove(dumpPath string) {
	utils.RemoveDirs([]string{dumpPath})
}

func (nw *Network) load(dumpPath string) error {
	nwConfigFile, err := os.Open(dumpPath)
	defer nwConfigFile.Close()
	if err != nil {
		log.Errorw(err, "Error to open network file", "path", dumpPath)
		return err
	}

	nwJson := make([]byte, 2000)
	n, err := nwConfigFile.Read(nwJson)
	if err != nil {
		log.Errorw(err, "Error to read network file", "path", dumpPath)
		return err
	}

	err = json.Unmarshal(nwJson[:n], nw)
	if err != nil {
		log.Errorw(err, "Error to unmarshal network file", "path", dumpPath)
		return err
	}
	return nil
}
