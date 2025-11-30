package consts

const (
	PERM_0777   = 0777 // 用户具、组用户和其它用户都有读/写/执行权限
	PERM_0755   = 0755 // 用户具有读/写/执行权限，组用户和其它用户具有读写权限；
	PERM_0644   = 0644 // 用户具有读写权限，组用户和其它用户具有只读权限；
	PERM_0622   = 0622 // 用户具有读/写权限，组用户和其它用户具只写权限；
	PERM_0777_S = 0777 // 用户具、组用户和其它用户都有读/写/执行权限
	PERM_0755_S = 0755 // 用户具有读/写/执行权限，组用户和其它用户具有读写权限；
	PERM_0644_S = 0644 // 用户具有读写权限，组用户和其它用户具有只读权限；
	PERM_0622_S = 0622 // 用户具有读/写权限，组用户和其它用户具只写权限；
)

// container 相关的
const (
	INFO_LOCATION      = "/var/lib/docker-demo/containers/"
	INFO_LOCATION_TEMP = INFO_LOCATION + "%s/"
	BASH               = "/proc/self/exe"
	LOG_FILE_TEMP      = INFO_LOCATION_TEMP + "%s-json.log"
)

// rootfs 相关的
const (
	IMAGE_PATH             = "/var/lib/docker-demo/image/"
	ROOT_PATH              = "/var/lib/docker-demo/overlay2/"
	LOWER_DIR_TEMP         = ROOT_PATH + "%s/lower"
	UPPER_DIR_TEMP         = ROOT_PATH + "%s/upper"
	WORK_DIR_TEMP          = ROOT_PATH + "%s/work"
	MERGED_DIR_TEMP        = ROOT_PATH + "%s/merged"
	OVERLAY_PARAMETER_TEMP = "lowerdir=%s,upperdir=%s,workdir=%s"
)

const (
	FDINDEX = 3 // index 为三的文件描述符，传递进来管道的另一端，默认会包含标准输入，标准输出，标准错误
)
