package network

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/pachirode/pkg/log"
	"github.com/vishvananda/netlink"

	"github.com/pachirode/docker-demo/pkg/errors"
)

type BridgeNetworkDriver struct {
}

func (bn *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (bn *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	n := &Network{
		Name:    name,
		IpRange: ipRange,
	}
	err := initBridge(n)
	if err != nil {
		log.Errorw(err, "Error to init bridge", "network", n)
	}
	return n, err
}

func initBridge(n *Network) error {
	// 创建虚拟设备
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return errors.WithMessage(err, "Error to add bridge")
	}

	// 设置虚拟设备地址和路由
	gatewayIP := *n.IpRange
	gatewayIP.IP = n.IpRange.IP
	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return errors.WithMessage(err, "Error set interface ip")
	}

	// 启动虚拟设备
	if err := setInterfaceUp(bridgeName); err != nil {
		return errors.WithMessage(err, "Error set interface up")
	}
	// 配置 iptables
	if err := setupIPTables(bridgeName, n.IpRange); err != nil {
		return errors.WithMessage(err, "Error setup iptables")
	}

	return nil
}

func createBridgeInterface(bridgeName string) error {
	// 存在同名的返回创建错误
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return errors.WithMessage(err, "Error to create, interface exists")
	}

	la := netlink.NewLinkAttrs()
	la.Name = bridgeName

	// 创建虚拟网络设备，相当于使用 ip link add XXX
	br := &netlink.Bridge{LinkAttrs: la}
	if err = netlink.LinkAdd(br); err != nil {
		return errors.WithMessage(err, "Error to create bridge")
	}
	return nil
}

func setInterfaceIP(bridgeName, rawIp string) error {
	retries := 2
	var inter netlink.Link
	var err error

	for i := 0; i < retries; i++ {
		inter, err = netlink.LinkByName(bridgeName)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return errors.WithMessage(err, "Error to create bridge link from netlink, run [ip link] to troubleshoot")
	}

	// 给网络接口配置地址，相当于使用 ip addr add xxx
	ipNet, err := netlink.ParseIPNet(rawIp)
	if err != nil {
		return err
	}
	addr := &netlink.Addr{IPNet: ipNet}
	return netlink.AddrAdd(inter, addr)
}

func setInterfaceUp(interfaceName string) error {
	inter, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return errors.WithMessage(err, "Error to find interface")
	}

	// 相当于 ip link set xxx up
	if err = netlink.LinkSetUp(inter); err != nil {
		return errors.WithMessage(err, "Error to enable interface")
	}

	return nil
}

func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
