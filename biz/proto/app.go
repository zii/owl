package proto

// screenshot: /screen/snap.png
type Device struct {
	Serial        string `json:"serial"`   // 序列号: FJR6R19B23007240
	Name          string `json:"name"`     // 设备名称:AQM-AL00
	Model         string `json:"model"`    // 型号:AQM_AL00
	ScreenSize    string `json:"screen"`   // 屏幕宽高
	Battery       int    `json:"battery"`  // 电量: 95 电量95%
	Brand         string `json:"brand"`    // 品牌HUAWEI
	CpuInfo       string `json:"cpu_info"` // 处理器[HUAWEI Kirin 710F]
	Operator      string `json:"operator"` // 运营商[,中国移动]
	Mem           string `json:"mem"`      // 总内存
	Disk          string `json:"disk"`     // 总磁盘G
	DiskAvailable string `json:"diska"`    // 可用空间G
	IP            string `json:"ip"`
	Uptime        int    `json:"uptime"`   // 开机时间(秒)
	Snapshot      string `json:"snapshot"` // 屏幕截图URL
}

type LabelIcon struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
}
