package main

import (
	"fmt"
	"log"
	"strings"
)

type BrightnessCommand int

const (
	DecreaseBrightness BrightnessCommand = iota
	IncreaseBrightness
)

type Monitor interface {
	getBrightness() int
	setBrightness(int)
}

func getCursorMonitor() Monitor {
	hMonitor, monitorDisplayName, err := cursorOnMonitor()
	if err != nil {
		log.Fatal(err)
	}

	// Suffixing with `\\Monitor0` is needed otherwise it'd report some other device ID.
	// 0 and 1 are options provided by https://github.com/posthumz/DisplayDevices (see readme).
	display := getDisplayFromName(monitorDisplayName+"\\Monitor0", 0)
	if display == nil {
		log.Fatal(fmt.Errorf("couldn't retrieve display's ID"))
	}
	deviceInstanceID := strings.Split(display.DeviceIDString(), "\\")[1] // no idx checking was done here

	monitorInstanceName, err := getMonitorInstanceName(deviceInstanceID)
	if err != nil {
		log.Fatal(err)
	}

	isWMIMonitor, err := isTypeWMIMonitor(monitorInstanceName)
	if err != nil {
		log.Fatal(err)
	}

	if isWMIMonitor {
		return WMIMonitor(monitorInstanceName)
	} else {
		return DDCMonitor(*hMonitor)
	}
}
