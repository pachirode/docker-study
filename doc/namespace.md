目前在 `Linux` 内核中主要实现了八种命名空间

- `Mount namespaces`
- `UTS namespaces`
- `IPC namespaces`
- `PID namespaces`
- `Network namespaces`
- `User namespaces`
- `Cgroup namespace`
- `Time namespace`

和命名空间相关的 `API`

- `clone`
    - 创建一个新的进程并把他加入到新的 `namespace` 中，由 `flag` 指定需要创建哪些 `namespace`
- `setns`
    - 将当前进程加入已有的 `namespace`
- `unshare`
    - 将当前进程移动到新创建的 `namespace`，由 `flag` 指定需要创建的 `namespace`
- `ioctl_ns`
    - 查询 `namespace` 信息

### 查看命名空间

```bash
ls /proc/{pid}/ns -al # 查看进程命名空间的信息
total 0
dr-x--x--x 2 root root 0 Nov 24 10:49 .
dr-xr-xr-x 9 root root 0 Nov 19 08:57 ..
lrwxrwxrwx 1 root root 0 Nov 25 14:36 cgroup -> 'cgroup:[4026531835]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 ipc -> 'ipc:[4026531839]' # 如果两个进程数字相同说明他们属于同一个命名空间
lrwxrwxrwx 1 root root 0 Nov 25 14:36 mnt -> 'mnt:[4026531841]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 net -> 'net:[4026531840]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 pid -> 'pid:[4026531836]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 pid_for_children -> 'pid:[4026531836]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 time -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 time_for_children -> 'time:[4026531834]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 user -> 'user:[4026531837]'
lrwxrwxrwx 1 root root 0 Nov 25 14:36 uts -> 'uts:[4026531838]'
```

`cat /proc/sys/user/max_pid_namespaces` 定义了限制命名空间的数量