package config

type Info struct {
	Pid         string   `json:"pid"`          // 容器的init进程在宿主机上的 PID
	Id          string   `json:"id"`           // 容器Id
	Name        string   `json:"name"`         // 容器名
	Command     string   `json:"command"`      // 容器内init运行命令
	CreatedTime string   `json:"created_time"` // 创建时间
	Status      string   `json:"status"`       // 容器的状态
	Volume      string   `json:"volume"`       // 挂载目录
	Image       string   `json:"image"`        // ImageName
	NetworkName string   `json:"networkName"`  // 容器所在的网络
	PortMapping []string `json:"portmapping"`  // 端口映射
	IP          string   `json:"ip"`
}
