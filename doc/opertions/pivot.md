实现容器和文件系统隔离

`busybox` 是一个最精简的镜像

```bash
docker pull busybox
docker export -o busybox.tar busybox
tar -zxvf busybox.tar -C /root/busybox/
```

### pivot_root

一个系统调用，主要功能是去改变当前 `root` 文件系统
两个文件夹不能同时存在当前 `root` 同一个文件系统，会把整个系统都给切换，原来对于旧系统的依赖会全部移除

- `put_old` 文件夹
    - 存放当前进程的 `root` 文件系统
- `new_root` 文件夹
    - 新的 `root` 文件系统

### chroot

只针对某个进程进行操作，其他部分可以运行在旧的系统上

# overlay 或者 aufs

创建容器文件系统，实现容器和本地文件系统隔离

- 创建只读层
- 创建容器读写层
- 将两层挂载到一个挂载点上

### 挂载标签
将一个目录或者文件系统挂载到另一个目录的技术，可以用来持久化容器中的数据

```bash
# 挂载到另一个目录，两个目录数据将会同步
mount -o bind /source/directory/ /target/directory/
```

