# Docker demo

基于 Go 1.24.6 和 debian12 实现一版简易的 Docker 来学习 docker 的底层原理

参考项目[mydocker](https://github.com/lixd/mydocker)

### 前置内容

- [namespace](doc/namespace.md)
- [cgroup](doc/cgroup.md)
- [fs](doc/cgroup.md)
- [网络](doc/网络.md)

### 问题

##### fork 子进程

`bash init [args]`

把参数全部跟在 `init` 后面，作为 `init` 参数，如何在 `init` 命令中解析参数
如果用户输入参数特别长，某些特殊字符可能失效
使用管道来实现父进程和子进程之间数据的传递
[匿名管道](doc/匿名管道.md)
