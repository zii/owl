// An app demonstrating most of the library's features.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	adb "github.com/zach-klippenstein/goadb"
)

var (
	port = flag.Int("p", adb.AdbPort, "")

	client *adb.Adb
)

func main() {
	flag.Parse()

	var err error
	client, err = adb.NewWithConfig(adb.ServerConfig{
		Port: *port,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting server…")
	fmt.Println("port:", *port)
	client.StartServer()

	serverVersion, err := client.ServerVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server version:", serverVersion)

	devices, err := client.ListDevices()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Devices:")
	for _, device := range devices {
		fmt.Printf("\t%+v\n", *device)
	}

	PrintDeviceInfoAndError(adb.AnyDevice())
	PrintDeviceInfoAndError(adb.AnyLocalDevice())
	PrintDeviceInfoAndError(adb.AnyUsbDevice())

	serials, err := client.ListDeviceSerials()
	if err != nil {
		log.Fatal(err)
	}
	for _, serial := range serials {
		PrintDeviceInfoAndError(adb.DeviceWithSerial(serial))
	}

	fmt.Println()
	fmt.Println("Watching for device state changes.")
	watcher := client.NewDeviceWatcher()
	for event := range watcher.C() {
		fmt.Printf("\t[%s]%+v\n", time.Now(), event)
	}
	if watcher.Err() != nil {
		printErr(watcher.Err())
	}

	//fmt.Println("Killing server…")
	//client.KillServer()
}

func printErr(err error) {
	fmt.Println("error:", err)
}

func PrintDeviceInfoAndError(descriptor adb.DeviceDescriptor) {
	device := client.Device(descriptor)
	if err := PrintDeviceInfo(device); err != nil {
		log.Println(err)
	}
}

func PrintDeviceInfo(device *adb.Device) error {
	serialNo, err := device.Serial()
	if err != nil {
		return err
	}
	devPath, err := device.DevicePath()
	if err != nil {
		return err
	}
	state, err := device.State()
	if err != nil {
		return err
	}

	fmt.Println(device)
	fmt.Printf("\tserial no: %s\n", serialNo)
	fmt.Printf("\tdevPath: %s\n", devPath)
	fmt.Printf("\tstate: %s\n", state)

	cmdOutput, err := device.RunCommand("pwd")
	if err != nil {
		fmt.Println("\terror running command:", err)
	}
	fmt.Printf("\tcmd output: %s\n", cmdOutput)

	stat, err := device.Stat("/sdcard")
	if err != nil {
		fmt.Println("\terror stating /sdcard:", err)
	}
	fmt.Printf("\tstat \"/sdcard\": %+v\n", stat)

	fmt.Println("\tfiles in \"/\":")
	entries, err := device.ListDirEntries("/")
	if err != nil {
		fmt.Println("\terror listing files:", err)
	} else {
		for entries.Next() {
			fmt.Printf("\t%+v\n", *entries.Entry())
		}
		if entries.Err() != nil {
			fmt.Println("\terror listing files:", err)
		}
	}

	fmt.Println("\tnon-existent file:")
	stat, err = device.Stat("/supercalifragilisticexpialidocious")
	if err != nil {
		fmt.Println("\terror:", err)
	} else {
		fmt.Printf("\tstat: %+v\n", stat)
	}

	fmt.Print("\tload avg: ")
	loadavgReader, err := device.OpenRead("/proc/loadavg")
	if err != nil {
		fmt.Println("\terror opening file:", err)
	} else {
		loadAvg, err := ioutil.ReadAll(loadavgReader)
		if err != nil {
			fmt.Println("\terror reading file:", err)
		} else {
			fmt.Println(string(loadAvg))
		}
	}

	return nil
}
