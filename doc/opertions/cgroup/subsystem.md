# Subsystem 接口

将 `cgroup` 抽象成 `path`，虚拟文件系统的路径便是 `cgroup` 在节点中的路径

### GetCgroupPath

获取当前 `subsystem` 在虚拟文件系统中的路径

- 找到对应 `subsystem` 挂载的 `hierarchy` 相对路径对应的虚拟文件系统中的地址
    - `/proc/self/mountinfo`
- 可以找出与当前进程相关的 `mount`
- 通过这个目录去操作 `cgroup`
