`UST Namespace` 是用来隔离 `nodename` 和 `domainname` 两个系统标识的

使用 `root` 权限运行代码，查看是否进入新的命名空间

```bash
pstree -pl | grep su
           |           `-sshd(81339)---sshd(81354)---bash(81357)---su(86911)---bash(86912)---go(86916)-+-uts(86966)-+-bash(86972)-+-grep(86976)

# 查看进程数值不同，因此不在一个命名空间中
readlink /proc/86972/ns/uts
uts:[4026532952]
readlink /proc/86966/ns/uts
uts:[4026531838]

hostname test # 测试外层 hostname 是否被修改
```

`IPC Namespace` 隔离 `Sys V IPC` 和 `POSIX message queues`

```bash
ipcs -q # 查看
ipcmk -Q # 创建
```

