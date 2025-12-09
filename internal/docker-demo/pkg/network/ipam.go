package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pachirode/pkg/log"

	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/errors"
)

type IIPAM interface {
	Allocate(subnet *net.IPNet) (ip net.IP, err error) // 从指定的 subnet 网段中分配 IP 地址
	Release(subnet *net.IPNet, ipaddr *net.IP) error   // 从指定的 subnet 网段中释放掉指定的 IP 地址
}

type IPAM struct {
	SubnetAllocatorPath string             // 分配文件存放位置
	Subnets             *map[string]string // 网段和位图算法的数组 map, key 是网段， value 是分配的位图数组
}

func (ipam *IPAM) load() error {
	// 文件不存在不需要加载
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		log.Errorw(err, "Error to open subnet config file", "path", ipam.SubnetAllocatorPath)
		return err
	}
	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		log.Errorw(err, "Error to read subnet config file", "path", ipam.SubnetAllocatorPath)
		return err
	}

	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Errorw(err, "Error dump allocation info", "data", subnetJson[:n])
		return err
	}
	return nil
}

func (ipam *IPAM) dump() error {
	utils.MkdirAll(consts.NETWORK_IPAM, consts.PERM_0644)

	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, consts.PERM_0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}

	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}

	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}

	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (ip net.IP, err error) {
	ipam.Subnets = &map[string]string{}

	err = ipam.load()
	if err != nil {
		log.Errorw(err, "Error to dump allocation info")
	}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	// 返回子网掩码255对应的位数和总位数
	one, size := subnet.Mask.Size()
	// 如果之前没分配过该网段，初始化全部网段
	if _, exist := (*ipam.Subnets)[subnet.String()]; !exist {
		// 用“0”填满这个网段的配置，uint8(size - one ）表示这个网段中有多少个可用地址
		// size - one是子网掩码后面的网络位数，2^(size - one)表示网段中的可用IP数
		// 而2^(size - one)等价于1 << uint8(size - one)
		// 左移一位就是扩大两倍
		ipCount := 1 << uint8(size-one)
		ipalloc := strings.Repeat("0", ipCount-2) // 减去 0和 255这俩个不可分配地址，所有可分配的 IP 地址
		// 初始化分配配置，标记 .0 和 .255 位置为不可分配，直接置为 1
		(*ipam.Subnets)[subnet.String()] = fmt.Sprintf("1%s1", ipalloc)
	}

	// 遍历位图数组，1 表示已经被分配，0 表示未分配
	for c := range (*ipam.Subnets)[subnet.String()] {
		if (*ipam.Subnets)[subnet.String()][c] == '0' {
			ipalloc := []byte((*ipam.Subnets)[subnet.String()])
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc) // 字符串无法修改，需要先转化为 byte 数组
			ip = subnet.IP                                     // 该网段的初始 IP
			for t := uint(4); t > 0; t -= 1 {
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			break
		}
	}

	if ip == nil {
		return nil, errors.WithMessage(err, "Error no available ip in subnet")
	}

	if err = ipam.dump(); err != nil {
		log.Errorw(err, "Error to dump ipam")
	}
	return
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		return errors.WithMessage(err, "Error to load ipam")
	}

	c := 0
	releaseIP := ipaddr.To4()
	for t := uint(4); t > 0; t -= 1 {
		c += int(releaseIP[t-1]-subnet.IP[t-1]) << ((4 - t) * 8)
	}

	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	if err = ipam.dump(); err != nil {
		return errors.WithMessage(err, "Error to dump ipam")
	}
	return nil
}
