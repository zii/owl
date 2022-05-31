package android

import (
	"embed"
	"io"
	"log"
	"strings"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

//go:embed busybox
var f_busybox embed.FS

const busyboxPath = "/data/local/tmp/busybox"

func FileExists(device *adb.Device, path string) bool {
	s, err := device.RunCommand("ls", path)
	if err != nil {
		log.Println("FileExists:", err)
		return false
	}
	if strings.Index(s, "No such file") >= 0 {
		return false
	}
	return true
}

func InstallEmbed(device *adb.Device, fs embed.FS, filename string, to string) bool {
	if FileExists(device, to) {
		return true
	}
	w, err := device.OpenWrite(to, 0755, time.Now())
	if err != nil {
		log.Println("InstallEmbed:", err)
		return false
	}
	defer w.Close()
	f, err := fs.Open(filename)
	if err != nil {
		log.Println("embed.Open:", err)
		return false
	}
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println("io.Copy:", err)
		return false
	}
	log.Printf("install %s success!\n", filename)
	return true
}

func InstallBusybox(device *adb.Device) bool {
	return InstallEmbed(device, f_busybox, "busybox", busyboxPath)
}
