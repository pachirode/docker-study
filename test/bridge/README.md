根据网络拓扑实现网络模型

### veth-pair

一对虚拟设备接口，成对出现，用于跨命名空间通信

##### 测试

```bash
# 创建命名空间
ip netns add netns-test
ip netns list
# 创建虚拟设备对
ip link add veth0 type veth peer name veth1
ip link show
# 切换命名空间
ip link set veth1 netns netns
# 命名空间中查看
ip netns exec netns-test ip link show

# 分配 IP
ip netns exec netns-test ip addr add 10.1.1.1/24 dev veth1
ip addr add 10.1.1.2/24 dev veth0
ip netns exec netns-test ip link set dev veth1 up
ip link set dev veth0 up

# 查看虚拟设备另一端
ip netns exec netns-test ethtool -S veth1
ip link | grep 128
```

### bridge

`sudo apt-get install bridge-utils`

```bash
# 创建网桥
sudo brctl addbr br-test
# 将 Veth 的另一端接入网桥
sudo brctl addif br-test veth1
```

### NAT

网络地址转换，容器内的 `IP` 和宿主机 `IP` 不一致，需要对 `IP` 层的源和目标 `IP` 进行转换

```bash
# 让非 172.18.0.0/24 网段的数据包路由给网桥，数据可以来到宿主机上
sudo ip netns exec netns-test ip route add default via 172.18.0.1 dev veth0

# 内核允许开启转发功能
sudo sysctl net.ipv4.conf.all.forwarding=1
```