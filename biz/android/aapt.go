package android

import (
	"embed"
	"fmt"
	"io"
	"log"
	"path"
	"strings"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

//go:embed aapt-arm-pie
var aapt_arm embed.FS

const aaptPath = "/data/local/tmp/aapt"

func InstallAAPT(device *adb.Device) bool {
	w, err := device.OpenWrite(aaptPath, 0755, time.Now())
	if err != nil {
		log.Println("InstallAAPT:", err)
		return false
	}
	defer w.Close()
	f, err := aapt_arm.Open("aapt-arm-pie")
	if err != nil {
		log.Println("aapt_arm.Open:", err)
		return false
	}
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println("io.Copy:", err)
		return false
	}
	log.Println("install aapt success!")
	return true
}

func AAPTExists(device *adb.Device) bool {
	s, err := device.RunCommand("ls", aaptPath)
	if err != nil {
		log.Println("AAPTExists:", err)
		return false
	}
	if strings.Index(s, "No such file") >= 0 {
		return false
	}
	return true
}

func CheckAAPT(device *adb.Device) bool {
	if AAPTExists(device) {
		return true
	}
	return InstallAAPT(device)
}

func ReadLabelIcon(device *adb.Device, apk string) (string, string) {
	s, err := device.RunCommand(aaptPath, "d", "badging", apk)
	if err != nil {
		log.Println("ReadLabelIcon:", err)
		return "", ""
	}
	i := strings.Index(s, "application: label='")
	if i < 0 {
		return "", ""
	}
	s = s[i+20:]
	i = strings.Index(s, "'")
	if i < 0 {
		return "", ""
	}
	label := s[:i]
	s = s[i:]
	i = strings.Index(s, "icon='")
	if i < 0 {
		return "", ""
	}
	s = s[i+6:]
	i = strings.Index(s, "'")
	if i < 0 {
		return "", ""
	}
	icon := s[:i]
	return label, icon
}

func extractIcon(device *adb.Device, apk string, icon string) []byte {
	s, err := device.RunCommand(busyboxPath, "unzip", apk, icon, "-p")
	if err != nil {
		log.Println("extractIcon:", adb.ErrorWithCauseChain(err))
		return nil
	}
	return []byte(s)
}

// 保存图标到 /screen/<pname>.ext
func SaveIcon(device *adb.Device, pname, apk string, icon string) string {
	data := extractIcon(device, apk, icon)
	if len(data) == 0 {
		return ""
	}
	ext := path.Ext(icon)
	filename := fmt.Sprintf("screen/%s%s", pname, ext)
	err := mfs.WriteFile(filename, data, 0755)
	if err != nil {
		log.Println("mfs.WriteFile icon:", err)
		return ""
	}
	return filename
}

func TestLabelIcon() {
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	apk := "/data/app/com.xinhang.mobileclient-90f08RnAdHpNwxzEJsefIg==/base.apk"
	label, icon := ReadLabelIcon(device, apk)
	log.Println("test label:", label, icon)
	filename := SaveIcon(device, "hehe", apk, icon)
	log.Println("icon filename:", filename)
}
