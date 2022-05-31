package android

import (
	"log"
	"strconv"
	"strings"
	"sync"

	adb "github.com/zach-klippenstein/goadb"
)

type PackageRaw struct {
	Name         string `json:"name"`  // 包名 com.huawei.appmarket
	Label        string `json:"label"` // 应用名
	Icon         string `json:"icon"`  // 图标url
	path         string // apk路径
	Size         string `json:"size"`           // 包大小
	InstallTime  string `json:"install_time"`   // 安装时间
	SdkVersion   string `json:"sdk_version"`    // SDK版本 minSdk=21 targetSdk=28
	Version      string `json:"version"`        // 版本 10.8.18
	TimeUsed     string `json:"time_used"`      // 使用时长
	LastTimeUsed string `json:"last_time_used"` // 上次使用时间
	Launch       int    `json:"launch"`         // 启动次数
}

var pkgcache map[string]*PackageRaw

func init() {
	pkgcache = make(map[string]*PackageRaw)
}

// read: SdkVersion, Version, InstallTime
func DumpPackages(device *adb.Device) map[string]*PackageRaw {
	var out = make(map[string]*PackageRaw)
	s, err := device.RunCommand("dumpsys", "package")
	if err != nil {
		log.Println("RunCommand dumpsys package:", err)
		return out
	}
	for {
		var p = &PackageRaw{}
		i := strings.Index(s, "  Package [")
		if i < 0 {
			break
		}
		s = s[i+11:]
		i = strings.Index(s, "]")
		if i < 0 {
			break
		}
		p.Name = s[:i]
		s = s[i:]
		//
		i = strings.Index(s, "versionCode=")
		if i < 0 {
			break
		}
		s = s[i+12:]
		i = strings.Index(s, "\n")
		if i < 0 {
			break
		}
		sdk := strings.TrimSpace(s[:i])
		if len(sdk) > 2 {
			ss := strings.SplitN(sdk, " ", 2)
			if len(ss) == 2 {
				sdk = ss[1]
			}
		}
		p.SdkVersion = sdk
		//
		i = strings.Index(s, "versionName=")
		if i < 0 {
			break
		}
		s = s[i+12:]
		i = strings.Index(s, "\n")
		if i < 0 {
			break
		}
		p.Version = strings.TrimSpace(s[:i])
		//
		i = strings.Index(s, "firstInstallTime=")
		if i < 0 {
			break
		}
		s = s[i+17:]
		i = strings.Index(s, "\n")
		if i < 0 {
			break
		}
		p.InstallTime = strings.TrimSpace(s[:i])
		i = strings.Index(s, "firstInstallTime=")
		if i < 0 {
			break
		}

		out[p.Name] = p
	}
	return out
}

// read: TimeUsed, LastTimeUsed, Launch
func DumpUsagestats(device *adb.Device) map[string]*PackageRaw {
	var out = make(map[string]*PackageRaw)
	s, err := device.RunCommand("dumpsys", "usagestats")
	if err != nil {
		log.Println("RunCommand dumpsys usagestats:", err)
		return out
	}
	i := strings.Index(s, "In-memory yearly stats")
	if i < 0 {
		return out
	}
	s = s[i+22:]
	i = strings.Index(s, "ChooserCounts")
	if i < 0 {
		i = strings.Index(s, "Settings:")
		if i < 0 {
			return out
		}
	}
	s = s[:i]
	for {
		p := &PackageRaw{}
		i := strings.Index(s, "package=")
		if i < 0 {
			break
		}
		s = s[i+8:]
		i = strings.Index(s, " ")
		if i < 0 {
			break
		}
		p.Name = s[:i]
		//
		i = strings.Index(s, "totalTimeUsed=\"")
		if i < 0 {
			i = strings.Index(s, "totalTime=\"")
			if i < 0 {
				break
			} else {
				s = s[i+11:]
			}
		} else {
			s = s[i+15:]
		}

		i = strings.Index(s, "\"")
		if i < 0 {
			break
		}
		p.TimeUsed = s[:i]
		//
		i = strings.Index(s, "lastTimeUsed=\"")
		if i < 0 {
			i = strings.Index(s, "lastTime=\"")
			if i < 0 {
				break
			} else {
				s = s[i+10:]
			}
		} else {
			s = s[i+14:]
		}
		i = strings.Index(s, "\"")
		if i < 0 {
			break
		}
		p.LastTimeUsed = s[:i]
		//
		i = strings.Index(s, "appLaunchCount=")
		if i > 0 {
			s = s[i+15:]
			i = strings.IndexAny(s, " \n")
			if i < 0 {
				break
			}
			p.Launch, _ = strconv.Atoi(s[:i])
		}
		out[p.Name] = p
	}
	return out
}

func ReadFileSize(device *adb.Device, path string) string {
	s, err := device.RunCommand("du", "-h", path)
	if err != nil {
		log.Println("RunCommand du -h:", path, err)
		return ""
	}
	fields := strings.Fields(s)
	if len(fields) != 2 {
		return ""
	}
	return fields[0]
}

func ReadPakcageNames(device *adb.Device) []*PackageRaw {
	var out = []*PackageRaw{}
	s, err := device.RunCommand("pm", "list", "packages", "-u", "-3", "-f")
	if err != nil {
		log.Println("RunCommand pm list packages -u -3 -f:", err)
		return out
	}
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.TrimPrefix(l, "package:")
		i := strings.LastIndex(l, "=")
		if i < 0 {
			break
		}
		path := l[:i]
		name := l[i+1:]
		p := &PackageRaw{
			Name: name,
			path: path,
		}
		out = append(out, p)
	}
	return out
}

func ReadAllPackages(device *adb.Device) map[string]struct{} {
	var out = make(map[string]struct{})
	s, err := device.RunCommand("pm", "list", "packages", "-u")
	if err != nil {
		log.Println("RunCommand pm list packages -u -s:", err)
		return out
	}
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.TrimPrefix(l, "package:")
		out[l] = struct{}{}
	}
	return out
}

// 获取已卸载的包
func ReadUninstalledPackages(device *adb.Device) []string {
	var out []string
	installed := ReadAllPackages(device)
	dups := make(map[string]struct{})
	s, err := device.RunCommand("dumpsys", "package")
	if err != nil {
		log.Println("RunCommand dumpsys package:", err)
		return out
	}
	i := strings.Index(s, "Package Changes:")
	if i < 0 {
		return out
	}
	s = s[i+16:]
	i = strings.Index(s, "Frozen packages:")
	if i < 0 {
		return out
	}
	s = s[:i]
	for {
		i := strings.Index(s, "package=")
		if i < 0 {
			break
		}
		s = s[i+8:]
		i = strings.Index(s, "\n")
		if i < 0 {
			break
		}
		name := strings.TrimSpace(s[:i])
		if _, ok := dups[name]; ok {
			continue
		}
		if _, ok := installed[name]; ok {
			continue
		}
		out = append(out, name)
	}
	return out
}

type TLPlist struct {
	Plist     []*PackageRaw `json:"install"`
	Uninstall []string      `json:"uninstall"`
}

func GetTLPlist() (*TLPlist, error) {
	var out = &TLPlist{}
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	// install aapt
	aapt_ok := CheckAAPT(device)
	_ = aapt_ok
	InstallBusybox(device)
	plist := ReadPakcageNames(device)
	for _, p := range plist {
		c := pkgcache[p.Name]
		if c != nil && c.Label != "" {
			p.Label = c.Label
			p.Icon = c.Icon
		}
	}
	dp := DumpPackages(device)
	du := DumpUsagestats(device)
	for _, p := range plist {
		p2 := dp[p.Name]
		if p2 != nil {
			p.SdkVersion = p2.SdkVersion
			p.Version = p2.Version
			p.InstallTime = p2.InstallTime
		}
		u := du[p.Name]
		if u != nil {
			p.TimeUsed = u.TimeUsed
			p.LastTimeUsed = u.LastTimeUsed
			p.Launch = u.Launch
		}
	}
	szmap := ReadSizeParallel(device, plist, 20)
	for _, p := range plist {
		p.Size = szmap[p.Name]
	}
	for _, p := range plist {
		pkgcache[p.Name] = p
	}
	//ReadLabelParallel(device, plist[:10], 20)
	out.Plist = plist
	//o := DumpPackages(device)
	//o := DumpUsagestats(device)
	//o := ReadFileSize(device, "/data/app/com.cmbc.cc.mbank--ggnamIsmPnM36yWiySYNA==/base.apk")
	//o := PakcageNames(device)
	//o := ReadSystemPackages(device)
	//o := ReadUninstalledPackages(device)
	//log.Println("o:", o)
	out.Uninstall = ReadUninstalledPackages(device)

	return out, nil
}

func ReadSizeParallel(device *adb.Device, plist []*PackageRaw, co int) map[string]string {
	out := make(map[string]string)
	mu := sync.RWMutex{}
	ch := make(chan *PackageRaw, 10)
	wg := &sync.WaitGroup{}
	wg.Add(co)
	for i := 0; i < co; i++ {
		go func() {
			defer wg.Done()
			for p := range ch {
				size := ReadFileSize(device, p.path)
				mu.Lock()
				out[p.Name] = size
				mu.Unlock()
			}
		}()
	}
	for _, p := range plist {
		ch <- p
	}
	close(ch)
	wg.Wait()
	return out
}

func (p *PackageRaw) ReadLabelIcon(device *adb.Device) {
	label, icon := ReadLabelIcon(device, p.path)
	p.Label = label
	p.Icon = SaveIcon(device, p.Name, p.path, icon)
}

func ReadLabelParallel(device *adb.Device, plist []*PackageRaw, co int) {
	ch := make(chan *PackageRaw, 20)
	wg := &sync.WaitGroup{}
	wg.Add(co)
	for i := 0; i < co; i++ {
		go func() {
			defer wg.Done()
			for p := range ch {
				p.ReadLabelIcon(device)
			}
		}()
	}
	for _, p := range plist {
		ch <- p
	}
	close(ch)
	wg.Wait()
}

func GetLabelIcon(pname string) (string, string) {
	p := pkgcache[pname]
	if p != nil && p.Label != "" {
		return p.Label, p.Icon
	}
	if p == nil {
		p = &PackageRaw{}
	}
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	p.ReadLabelIcon(device)
	return p.Label, p.Icon
}
