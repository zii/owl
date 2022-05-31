package android

import (
	"fmt"
	"owl/base"

	adb "github.com/zach-klippenstein/goadb"
)

func TestLS() {
	descriptor := adb.AnyUsbDevice()
	device := client.Device(descriptor)
	entries, err := device.ListDirEntries("/sdcard/Pictures")
	base.Raise(err)
	for entries.Next() {
		all, _ := entries.ReadAll()
		fmt.Println(base.JsonString(all))
	}
}
