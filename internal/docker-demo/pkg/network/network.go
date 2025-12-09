package network

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/pachirode/pkg/log"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"

	"github.com/pachirode/docker-demo/internal/docker-demo/options"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/config"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/consts"
	"github.com/pachirode/docker-demo/internal/docker-demo/pkg/utils"
	"github.com/pachirode/docker-demo/pkg/errors"
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
	utils.RemoveDirs([]string{path.Join(dumpPath, nw.Name)})
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

func loadNetwork() (map[string]*Network, error) {
	networks := map[string]*Network{}

	err := filepath.Walk(consts.NETWORK_ROOT, func(netPath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, netName := path.Split(netPath)
		network := &Network{
			Name: netName,
		}

		if err = network.load(netPath); err != nil {
			log.Errorw(err, "Error to load net", "path", netPath)
			return err
		}

		if network.IpRange != nil {
			networks[netName] = network
		}
		return nil
	})

	return networks, err
}

func CreateNetwork(opts *options.NetworkOptions, name string) error {
	_, cidr, _ := net.ParseCIDR(opts.Subnet)
	ip, err := ipAllocator.Allocate(cidr)
	if err != nil {
		return err
	}
	cidr.IP = ip

	nw, err := drivers[opts.Driver].Create(cidr.String(), name)
	if err != nil {
		log.Errorw(err, "Error to create driver")
		return err
	}

	return nw.dump(consts.NETWORK_ROOT)
}

func ListNetwork() {
	networks, err := loadNetwork()
	if err != nil {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "NAME\tIpRange\tDriver\n")
	for _, net := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			net.Name,
			net.IpRange.String(),
			net.Driver,
		)
	}
	if err = w.Flush(); err != nil {
		log.Errorw(err, "Error to flush network")
		return
	}
}

func DeleteNetwork(networkName string) error {
	networks, err := loadNetwork()
	if err != nil {
		return errors.WithMessage(err, "Error to load networks")
	}

	network, ok := networks[networkName]
	if !ok {
		return errors.WithMessage(err, "Error to find network")
	}

	if err = ipAllocator.Release(network.IpRange, &network.IpRange.IP); err != nil {
		return errors.WithMessage(err, "Error to release gateway ip")
	}

	if err = drivers[network.Driver].Delete(network); err != nil {
		return errors.WithMessage(err, "Error to delete driver")
	}

	network.remove(consts.NETWORK_ROOT)
	return nil
}

// Connect 连接容器到已经创建好的网络
func Connect(networkName string, info *config.Info) (net.IP, error) {
	networks, err := loadNetwork()
	if err != nil {
		return nil, errors.WithMessage(err, "Error to load network")
	}
	network, ok := networks[networkName]
	if !ok {
		return nil, errors.WithMessage(err, "Error to found network")
	}

	ip, err := ipAllocator.Allocate(network.IpRange)
	if err != nil {
		return ip, errors.WithMessage(err, "Error to allocate ip")
	}

	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", info.Id, networkName),
		IPAddress:   ip,
		Network:     network,
		PortMapping: info.PortMapping,
	}

	if err = drivers[network.Driver].Connect(network.Name, ep); err != nil {
		return ip, err
	}

	if err = configEndpointIpAndRoute(ep, info); err != nil {
		return ip, err
	}

	return ip, configPortMapping(ep, false)
}

func Disconnect(networkName string, info *config.Info) error {
	networks, err := loadNetwork()
	if err != nil {
		return errors.WithMessage(err, "Error to load network")
	}
	network, ok := networks[networkName]
	if !ok {
		return errors.WithMessage(err, "Error to found network")
	}

	drivers[network.Driver].Disconnect(fmt.Sprintf("%s-%s", info.Id, networkName))
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", info.Id, networkName),
		IPAddress:   net.ParseIP(info.IP),
		Network:     network,
		PortMapping: info.PortMapping,
	}
	return configPortMapping(ep, true)
}

// enterContainerNetNS 将容器网络端点加入容器网络中，切换当前线程进入到容器网络空间
func enterContainerNetNS(enLink *netlink.Link, info *config.Info) func() {
	// 打开文件操作 Net 命名空间
	f, err := os.OpenFile(fmt.Sprintf("/proc/%s/ns/net", info.Pid), os.O_RDONLY, 0)
	if err != nil {
		log.Errorw(err, "Error to get net namespace", "pid", info.Pid)
	}

	nsFD := f.Fd()
	// 锁定当前线程 Pid，否则 goroutine 可能被调度到其他线程，需要保证线程始终在这个命名空间
	runtime.LockOSThread()

	// 修改网络端点Veth的另外一端，将其移动到容器的Net Namespace 中
	if err = netlink.LinkSetNsFd(*enLink, int(nsFD)); err != nil {
		log.Errorw(err, "Error to set veth in net namespace")
	}

	// 获取当前网络命名空间
	origns, err := netns.Get()
	if err != nil {
		log.Errorw(err, "Error to get current net namespace")
	}

	// 将当前进程加入网络命名空间
	if err = netns.Set(netns.NsHandle(nsFD)); err != nil {
		log.Errorw(err, "Error to set net namespace")
	}

	// 恢复到原始的网络命名空间
	return func() {
		// 恢复到上面获取到的之前的 Net Namespace
		netns.Set(origns)
		origns.Close()
		// 取消对当附程序的线程锁定
		runtime.UnlockOSThread()
		f.Close()
	}
}

// configEndpointIpAndRoute 配置网络端点和路由
func configEndpointIpAndRoute(ep *Endpoint, info *config.Info) error {
	// 网络端点中 veth 的另一端
	peerLink, err := netlink.LinkByName(ep.Device.PeerName)
	if err != nil {
		return errors.WithMessage(err, "Error to found veth")
	}

	defer enterContainerNetNS(&peerLink, info)()

	interIP := *ep.Network.IpRange
	interIP.IP = ep.IPAddress

	if err = setInterfaceIP(ep.Device.PeerName, interIP.String()); err != nil {
		return errors.WithMessage(err, "Error to set veth")
	}

	if err = setInterfaceUp(ep.Device.PeerName); err != nil {
		return errors.WithMessage(err, "Error to set veth up")
	}

	// 网络命名空间中的默认本地网卡为 lo，默认是关闭状态
	if err = setInterfaceUp("lo"); err != nil {
		return errors.WithMessage(err, "Error to set lo up")
	}

	// 设置容器内的外部请求都通过容器内的 veth 端点访问
	_, cidr, _ := net.ParseCIDR("0.0.0.0/0")
	// 构建要添加的路由数据，包括网络设备、网关IP及目的网段
	// 相当于route add -net 0.0.0.0/0 gw (Bridge网桥地址) dev （容器内的Veth端点设备)

	defaultRoute := &netlink.Route{
		LinkIndex: peerLink.Attrs().Index,
		Gw:        ep.Network.IpRange.IP,
		Dst:       cidr,
	}

	if err = netlink.RouteAdd(defaultRoute); err != nil {
		return errors.WithMessage(err, "Error to route add veth")
	}

	return nil
}

func configPortMapping(ep *Endpoint, isDelete bool) error {
	action := "-A"
	if isDelete {
		action = "-D"
	}

	var err error
	// 遍历容器端口映射列表
	for _, pm := range ep.PortMapping {
		// 分割成宿主机的端口和容器的端口
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			log.Errorw(err, "Error split port")
			continue
		}
		// 由于iptables没有Go语言版本的实现，所以采用exec.Command的方式直接调用命令配置
		// 在iptables的PREROUTING中添加DNAT规则
		// 将宿主机的端口请求转发到容器的地址和端口上
		// iptables -t nat -A PREROUTING ! -i testbridge -p tcp -m tcp --dport 8080 -j DNAT --to-destination 10.0.0.4:80
		iptablesCmd := fmt.Sprintf("-t nat %s PREROUTING ! -i %s -p tcp -m tcp --dport %s -j DNAT --to-destination %s:%s",
			action, ep.Network.Name, portMapping[0], ep.IPAddress.String(), portMapping[1])
		cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
		log.Infow("配置端口映射 DNAT cmd:", "cmd", cmd.String())
		// 执行iptables命令,添加端口映射转发规则
		output, err := cmd.Output()
		if err != nil {
			log.Errorw(err, "iptables Output", "out", output)
			continue
		}
	}
	return err
}
