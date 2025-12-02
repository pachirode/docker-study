# 操作 Cgroup

### 查看 Cgroup 版本

```bash
mount | grep cgroup
cgroup on /sys/fs/cgroup type cgroup (rw,nosuid,nodev,noexec,relatime,cgroup2) 
cgroup2 on /sys/fs/cgroup2 type cgroup2 (rw,nosuid,nodev,noexec,relatime) # cgroup2
```

### 创建 Cgroup
```bash
# 创建并挂起一个 hierarchy
mkdir cgroup-test
# 挂载一个 hierarchy
sudo mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test
# 挂载之后可以看到目录下生成一些默认文件，这些文件就是根节点的配置项
ls cgroup-test/
cgroup.clone_children  cgroup.procs  cgroup.sane_behavior  notify_on_release  release_agent  tasks
# 扩展两个子 cgroup，在一个 cgroup 创建文件夹，会把这个文件夹标记为子 cgroup
```

### 移动或者添加进程

```bash
echo &&
# 将所在终端的进程移动到 cgroup-1
sudo sh -c "echo $$ >> tasks"
cat /proc/130930/cgroup 
1:name=cgroup-test:/cgroup-1 # 已经被添加进来了
0::/user.slice/user-1000.slice/session-625.scope

```
