# Run 函数

这个是父进程，当前进程执行的内容
- `/proc/self/exe` 自己调用自己，对创建出来的进程进行初始化
- `init` 为传递到本进程的第一个参数，项目中是调用 `initCommand` 去初始化进程的一些环境和资源
- `clone` 参数，`fork` 一个新的进程，并使用 `namespace` 进行隔离
