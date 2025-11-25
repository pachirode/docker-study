`UnionFS` 把其他文件系统联合到一个联合挂载点的文件系统服务
将多个目录叠加到一块，形成一个单一一致的文件系统

### 写时复制

对可修改资源实现高效复制的资源管理技术，如果一个资源是重复的，但是没有任何修改，这个时候不需要创建一个新的资源，这个资源可以被新旧实例共享
创建新资源发生在第一次写操作，通过这种资源共享的方式，可以显著减少修改资源复制带来的消耗

### 文件系统

`overlayfs` 是一种类似 `aufs` 的堆叠文件系统，它依赖并建立在其他的文件系统上，并不直接参与磁盘空间结构的划分，仅仅将原来底层文件系统中的不同目录进行合并

- 上下层同名目录合并
- 上下层同名文件覆盖
- `lower dir` 文件写时拷贝

### overlayfs 测试

```bash
├── lower
│   ├── a
│   └── c
├── merged
├── upper
│   ├── a
│   └── b
└── work

sudo mount \
            -t overlay \ # 表示文件系统
            overlay \
            -o lowerdir=./lower,upperdir=./upper,workdir=./work \ # 指定 lowerdir、upperdir、workdir
            ./merged
```

- `lower dir`
  - 底层目录，只提供数据，不能写
  - 如果修改底层映射的文件，会复制一份到上层目录
- `upper dir`
  - 上层目录，可读可写

