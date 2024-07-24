package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// can take absolute value like "44" or relative like "-12" or "+13"
// for absolute: setBrightness(strconv.Itoa(brightness))
func setBrightness(value string) {
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

	cmd := exec.Command("monitorian", "/set", monitorInstanceName, value)
	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}
}
