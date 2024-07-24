package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func getCursorMonitor() string {
	monitorDisplayName, err := cursorOnMonitor()
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

	return monitorInstanceName
}

// can take absolute value like "44" or relative like "-12" or "+13"
// for absolute: setBrightness(strconv.Itoa(brightness))
func setBrightness(value string) int {
	monitorInstanceName := getCursorMonitor()

	cmd := exec.Command("monitorian", "/set", monitorInstanceName, value)
	outB, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	out := strings.TrimSpace(string(outB))
	return getBrightnessFromMonitorianOutput(out)
}

func getBrightnessFromMonitorianOutput(out string) int {
	/*
		the value can look like:
		DISPLAY\GSM5BC4\5&1dbfb976&0&UID4352 "LG IPS QHD" 10 B *
		or
		DISPLAY\GSM5BC4\5&1dbfb976&0&UID4352 "LG IPS QHD" 10 B
	*/
	numStartIdx := strings.LastIndex(out, "\"") + 2
	numEndIdx := strings.LastIndex(out, "B") - 1
	numStr := out[numStartIdx:numEndIdx]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		log.Fatal(err)
	}
	return clamp(0, 100, num) // we'd parse and clamp it because monitorian sometimes output -1 as brightness rather than 0
}
