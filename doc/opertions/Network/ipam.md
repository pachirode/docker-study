# bitmap 算法

位图算法，用于大规模连续且少状态的数据处理
用于 `IP` 地址分配中，状态有两种 `1` 已经分配，`0` 未被分配

通过遍历数组计算偏移量可以快速定位到目标 `IP` 地址

# 测试

```bash
# 查看创建的 bridge 设备
ip link show dev test-network
11: test-network: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN mode DEFAULT group default qlen 1000
    link/ether 22:39:c8:0b:75:f1 brd ff:ff:ff:ff:ff:ff

# 查看地址配置和路由配置
ip addr show dev test-network
11: test-network: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether 22:39:c8:0b:75:f1 brd ff:ff:ff:ff:ff:ff
    inet 192.168.10.4/24 brd 192.168.10.255 scope global test-network
       valid_lft forever preferred_lft forever
    inet 192.168.10.5/24 brd 192.168.10.255 scope global secondary test-network
       valid_lft forever preferred_lft forever
    inet 192.168.10.6/24 brd 192.168.10.255 scope global secondary test-network
       valid_lft forever preferred_lft forever
    inet6 fe80::8c88:56ff:fe8b:2eda/64 scope link
       valid_lft forever preferred_lft forever

# 查看 iptables
iptables -t nat -vnL POSTROUTING
Chain POSTROUTING (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 MASQUERADE  0    --  *      !docker0  172.17.0.0/16        0.0.0.0/0
    0     0 MASQUERADE  0    --  *      !br-57b0985ada23  172.21.0.0/16        0.0.0.0/0
    0     0 MASQUERADE  0    --  *      !br-4cf94eaff650  172.24.0.0/16        0.0.0.0/0
    0     0 MASQUERADE  0    --  *      !br-db426ba1c71e  172.22.0.0/16        0.0.0.0/0
    0     0 MASQUERADE  0    --  *      !br-cb4536b203d0  172.20.0.0/16        0.0.0.0/0
    0     0 MASQUERADE  0    --  *      !br-7cb90afafef3  172.19.0.0/16        0.0.0.0/0
```