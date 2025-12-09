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

// Name 返回网络驱动名
func (bn *BridgeNetworkDriver) Name() string {
	return "bridge"
}

// Create 创建网络
func (bn *BridgeNetworkDriver) Create(subnet, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip
	n := &Network{
		Name:    name,
		IpRange: ipRange,
		Driver:  bn.Name(),
	}
	err := bn.initBridge(n)
	if err != nil {
		log.Errorw(err, "Error to init bridge", "network", n)
	}
	return n, err
}

// Delete 删除网络
func (bn *BridgeNetworkDriver) Delete(network *Network) error {
	// 清除路由规则
	err := deleteIpRoute(network.Name, network.IpRange.String())
	if err != nil {
		return errors.WithMessage(err, "Error to del network route")
	}

	// 清除 iptables 规则
	err = setupIPTables(network.Name, network.IpRange, true)
	if err != nil {
		return errors.WithMessage(err, "Error to delete iptables rule")
	}

	err = bn.deleteBridge(network)
	if err != nil {
		return errors.WithMessage(err, "Error to delete bridge")
	}

	return nil
}

// Connect 将指定的 Endpoint 连接到指定的网络
func (bn *BridgeNetworkDriver) Connect(networkName string, endpoint *Endpoint) error {
	bridgeName := networkName
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	// 创建 Veth 接口配置
	la := netlink.NewLinkAttrs()
	// 接口名最大长度不能超过 15
	la.Name = endpoint.ID[:5]
	// 将 Veth 的一端挂载到网络对应的 Linux Bridge
	la.MasterIndex = br.Attrs().Index

	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  "cfi-" + endpoint.ID[:5],
	}

	// 创建 Veth 接口，因为指定 MasterIndex 使用已经将其挂载到网桥上
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return errors.WithMessage(err, "Error to add endpoint device")
	}

	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return errors.WithMessage(err, "Error to set endpoint device up")
	}
	return nil
}

// Disconnect 将 veth 从设备上解绑
func (bn *BridgeNetworkDriver) Disconnect(endpointID string) error {
	vethName := endpointID[:5]
	veth, err := netlink.LinkByName(vethName)
	if err != nil {
		return errors.WithMessage(err, "Error to find veth")
	}

	// 从网桥解绑
	err = netlink.LinkSetNoMaster(veth)
	if err != nil {
		return errors.WithMessage(err, "Error set veth no master")
	}

	// 删除 veth-pair，xxx 和 cif-xxx
	err = netlink.LinkDel(veth)
	if err != nil {
		return errors.WithMessage(err, "Error to del veth")
	}
	// 经过测试，另一个 veth 会被同步删除
	//veth2Name := "cfi-" + vethName
	//veth2, err := netlink.LinkByName(veth2Name)
	//if err != nil {
	//	return errors.WithMessage(err, "Error to find veth2")
	//}
	//err = netlink.LinkDel(veth2)
	//if err != nil {
	//	return errors.WithMessage(err, "Error to del veth2")
	//}

	return nil
}

func (bn *BridgeNetworkDriver) initBridge(n *Network) error {
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
	if err := setupIPTables(bridgeName, n.IpRange, false); err != nil {
		return errors.WithMessage(err, "Error setup iptables")
	}

	return nil
}

// deleteBridge 删除虚拟网桥设备
func (bn *BridgeNetworkDriver) deleteBridge(n *Network) error {
	bridgeName := n.Name

	la, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return errors.WithMessage(err, "Error to find link")
	}

	if err = netlink.LinkDel(la); err != nil {
		return errors.WithMessage(err, "Error to remove bridge interface")
	}

	return nil
}

// createBridgeInterface 创建 Bridge 设备
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

// setInterfaceIP 设置 Bridge 设备地址和路由
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
		return errors.WithMessage(err, "Error to find the bridge link from netlink, run [ip link] to troubleshoot")
	}

	ipNet, err := netlink.ParseIPNet(rawIp)
	if err != nil {
		return err
	}
	// 相当于 ip addr add xxx
	// 配置所在的网段信息 192.168.0.0/24
	// 配置路由表 192.168.0.0/24 转发到当前网络接口
	addr := &netlink.Addr{IPNet: ipNet}
	return netlink.AddrAdd(inter, addr)
}

// 删除路由规则 ip addr del xxx
func deleteIpRoute(name, rawIP string) error {
	reties := 2
	var inter netlink.Link
	var err error

	for i := 0; i < reties; i++ {
		inter, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return errors.WithMessage(err, "Error to find the bridge link from netlink, run [ip link] to troubleshoot")
	}

	// 删除所有可以匹配的设备
	list, err := netlink.RouteList(inter, netlink.FAMILY_V4)
	if err != nil {
		log.Errorw(err, "Error to route links")
		return err
	}

	for _, route := range list {
		if route.Dst.String() == rawIP { // 根据子网进行匹配
			err = netlink.RouteDel(&route)
			if err != nil {
				log.Errorf("route [%v] del failed,detail:%v", route, err)
				continue
			}
		}
	}

	return nil
}

// setInterfaceUp 启动 Bridge 设备
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

// 设置 iptables 对应的 bridge 规则
// iptables -t nat -A POSTROUTING -s 172.18.0.0/24 -o eth0 -j MASQUERADE
// iptables -t nat -A POSTROUTING -s {subnet} -o {deviceName} -j MASQUERADE
func setupIPTables(bridgeName string, subnet *net.IPNet, isDelete bool) error {
	action := "-A"
	if isDelete {
		action = "-D"
	}
	iptablesCmd := fmt.Sprintf("-t nat %s POSTROUTING -s %s ! -o %s -j MASQUERADE", action, subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	log.Infow("config SNAT cmd", "cmd", cmd.String())
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output, %v", output)
	}
	return err
}
