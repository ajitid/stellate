package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type MonitorInfo struct {
	name       string
	brightness int
}

type BrightnessCommand int

const (
	DecreaseBrightness BrightnessCommand = iota
	IncreaseBrightness
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

// Beside absolute value like "44", the command can also take relative like "-12" or "+13".
// This command can also return the output.
// This fn however, will only take absolute value
func setBrightness(monitorInstanceName string, value int) {
	cmd := exec.Command("monitorian", "/set", monitorInstanceName, strconv.Itoa(value))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func getBrightness(monitorInstanceName string) int {
	cmd := exec.Command("monitorian", "/get", monitorInstanceName)
	outB, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	out := strings.TrimSpace(string(outB))
	if out == "" {
		log.Fatal("brightness reported empty by monitorian")
	}
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
	return num
}

func brightnessSetter(commandChan <-chan BrightnessCommand) {
	var currentMonitor MonitorInfo
	resetTimer := make(chan bool)
	go clearCurrentMonitor(&currentMonitor, resetTimer)

	for {
		command := <-commandChan

		monitorInstanceName := getCursorMonitor()
		if currentMonitor.name != monitorInstanceName {
			b := getBrightness(monitorInstanceName)
			if b == -1 {
				log.Fatal("brightness reported -1 by monitorian")
			}

			resetTimer <- true
			currentMonitor.name = monitorInstanceName
			currentMonitor.brightness = b
		} else {
			resetTimer <- true
		}

		switch command {
		case DecreaseBrightness:
			currentMonitor.brightness =
				clamp(0, 100,
					int(math.Floor(
						snapNumber(6.25)(float64(currentMonitor.brightness)-6.25))))
			go setBrightness(currentMonitor.name, currentMonitor.brightness)
		case IncreaseBrightness:
			currentMonitor.brightness =
				clamp(0, 100,
					int(math.Floor(
						snapNumber(6.25)(float64(currentMonitor.brightness)+6.25))))
			go setBrightness(currentMonitor.name, currentMonitor.brightness)
		}
	}
}

func clearCurrentMonitor(currentMonitor *MonitorInfo, resetTimer <-chan bool) {
	t := time.AfterFunc(0, func() {})

	for {
		<-resetTimer
		t.Stop()
		t = time.AfterFunc(1*time.Second+200*time.Millisecond, func() {
			currentMonitor.name = ""
			currentMonitor.brightness = 0
		})
	}
}
