package android

import (
	"fmt"
	"log"
	"owl/biz/proto"
	"strconv"
	"strings"
	"sync"

	adb "github.com/zach-klippenstein/goadb"
)

var (
	client *adb.Adb
	tldev  proto.Device // cache
)

func AdbStart() error {
	port := adb.AdbPort
	c, err := adb.NewWithConfig(adb.ServerConfig{
		Port: port,
	})
	if err != nil {
		return err
	}
	err = c.StartServer()
	if err != nil {
		return err
	}
	client = c
	return nil
}

func GetScreenSize(device *adb.Device) string {
	s, err := device.RunCommand("wm", "size")
	if err != nil {
		log.Println("GetScreenSize err:", err)
		return ""
	}
	i := strings.Index(s, "size:")
	if i < 0 {
		return ""
	}
	s = strings.TrimSpace(s[i+5:])
	return s
}

func SaveScreenShot(device *adb.Device) {
	s, err := device.RunCommand("screencap", "-p")
	if err != nil {
		log.Println("err:", err)
	}
	data := []byte(s)

	filename := fmt.Sprintf("screen/snap.png")
	err = mfs.WriteFile(filename, data, 0755)
	if err != nil {
		log.Println("mfs.WriteFile:", err)
		return
	}
}

func GetBattery(device *adb.Device) int {
	s, err := device.RunCommand("dumpsys", "battery")
	if err != nil {
		log.Println("GetBattery err:", err)
		return 0
	}
	i := strings.Index(s, "level:")
	if i < 0 {
		return 0
	}
	s = s[i+6:]
	i = strings.Index(s, "\n")
	if i > 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	level, _ := strconv.Atoi(s)
	return level
}

func trimBrackets(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimLeft(s, "[")
	s = strings.TrimRight(s, "]")
	return s
}

func getprops(device *adb.Device) Props {
	out := make(Props)
	s, err := device.RunCommand("getprop")
	if err != nil {
		log.Println("getprops err:", err)
		return out
	}
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		p := strings.SplitN(l, ":", 2)
		if len(p) != 2 {
			continue
		}
		k := trimBrackets(p[0])
		v := trimBrackets(p[1])
		out[k] = PropValue(v)
	}
	return out
}

func GetMemory(device *adb.Device) string {
	s, _ := device.RunCommand("free")
	if s == "" {
		return ""
	}
	i := strings.Index(s, "Mem:")
	if i < 0 {
		return ""
	}
	s = s[i+4:]
	i = strings.Index(s, "\n")
	if i > 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	if len(fields) > 0 {
		f, _ := strconv.Atoi(fields[0])
		return humMemory(f)
	}
	return ""
}

func humMemory(s int) string {
	a := float64(s)
	if a < 1024 {
		return strconv.FormatFloat(a, 'f', -1, 64)
	}
	a /= 1024
	if a < 1024 {
		return strconv.FormatFloat(a, 'f', 1, 64) + "K"
	}
	a /= 1024
	if a < 1024 {
		return strconv.FormatFloat(a, 'f', 1, 64) + "M"
	}
	a /= 1024
	if a < 1024 {
		return strconv.FormatFloat(a, 'f', 1, 64) + "G"
	}
	a /= 1024
	return strconv.FormatFloat(a, 'f', 1, 64) + "T"
}

func humSizeK(s int) string {
	if s < 1024 {
		return fmt.Sprintf("%dK", s)
	}
	s /= 1024
	if s < 1024 {
		return fmt.Sprintf("%dM", s)
	}
	s /= 1024
	if s < 1024 {
		return fmt.Sprintf("%dG", s)
	}
	s /= 1024
	return fmt.Sprintf("%dT", s)
}

func GetDisk(device *adb.Device) (string, string) {
	s, _ := device.RunCommand("df", "/sdcard")
	if s == "" {
		return "", ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return "", ""
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return "", ""
	}
	f1, _ := strconv.Atoi(fields[1])
	f3, _ := strconv.Atoi(fields[3])
	return humSizeK(f1), humSizeK(f3)
}

func GetIP(device *adb.Device) string {
	s, _ := device.RunCommand("netcfg")
	if s == "" {
		return ""
	}
	i := strings.Index(s, "wlan0")
	if i < 0 {
		return ""
	}
	s = s[i+5:]
	i = strings.Index(s, "\n")
	if i > 0 {
		s = s[:i]
	}
	fields := strings.Fields(s)
	if len(fields) > 2 {
		return fields[1]
	}
	return ""
}

func GetUptime(device *adb.Device) int {
	s, _ := device.RunCommand("cat", "/proc/uptime")
	if s == "" {
		return 0
	}
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return 0
	}
	f, _ := strconv.ParseFloat(fields[0], 64)
	return int(f)
}

func GetTLDevice() (*proto.Device, error) {
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	info, err := device.DeviceInfo()
	if err != nil {
		log.Println("device info:", err)
		return nil, err
	}
	//if info.Serial == tldev.Serial {
	//	return &tldev, nil
	//}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		SaveScreenShot(device)
	}()
	out := &proto.Device{}
	out.Serial = info.Serial
	out.Model = info.Model
	out.Name = info.Product
	out.ScreenSize = GetScreenSize(device)
	out.Battery = GetBattery(device)
	props := getprops(device)
	out.Brand = props.Get("ro.product.brand").String()
	out.CpuInfo = props.Get("ro.config.cpu_info_display").String()
	if out.CpuInfo == "" {
		out.CpuInfo = props.Get("ro.hardware.alter").String()
	}
	out.Operator = props.Operator()
	out.Mem = GetMemory(device)
	out.Disk, out.DiskAvailable = GetDisk(device)
	out.IP = GetIP(device)
	out.Uptime = GetUptime(device)
	out.Snapshot = "/screen/snap.png"
	tldev = *out
	wg.Wait()
	return out, nil
}

func GetDeviceInfo() *adb.DeviceInfo {
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	info, err := device.DeviceInfo()
	if err != nil {
		return nil
	}
	return info
}

func NewWatcher() *adb.DeviceWatcher {
	watcher := client.NewDeviceWatcher()
	return watcher
}
