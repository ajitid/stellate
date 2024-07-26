package main

import (
	"fmt"
	"log"
	"math"
	"strings"
	"syscall"
	"time"

	retry "github.com/avast/retry-go/v4"
	"github.com/gek64/displayController"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type BrightnessCommand int

const (
	DecreaseBrightness BrightnessCommand = iota
	IncreaseBrightness
)

type Monitor interface {
	getBrightness() (int, error)
	setBrightness(int)
	getInstanceName() string
	getPosition() rl.Vector2
}

type MonitorInfo struct {
	name       string
	brightness int
	monitor    *Monitor
}

var currentMonitor MonitorInfo

func getCursorMonitor() Monitor {
	hMonitor, monitorDisplayName, pos, err := cursorOnMonitor()
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
		return WMIMonitor{name: monitorInstanceName, pos: pos}
	} else {
		physicalMonitor, err := displayController.GetPhysicalMonitor(syscall.Handle(*hMonitor))
		if err != nil {
			log.Fatal(err)
		}
		return DDCMonitor{name: monitorInstanceName, physicalMonitor: &physicalMonitor, pos: pos}
	}
}

func brightnessSetter(commandChan <-chan BrightnessCommand, popupVisibleChan chan<- bool, popupPosChan chan<- rl.Vector2) {
	resetTimer := make(chan bool)
	go clearCurrentMonitor(&currentMonitor, resetTimer, popupVisibleChan)

	for {
		command := <-commandChan

		m := getCursorMonitor()
		if currentMonitor.name != m.getInstanceName() {
			b, err := retry.DoWithData(
				func() (int, error) {
					return m.getBrightness()
				},
				retry.Attempts(uint(5)),
				retry.OnRetry(func(attempt uint, err error) {
					log.Println("retrying", attempt+1, "time to get brightness")
				}))
			if err != nil {
				log.Fatal(fmt.Errorf("getting monitor brightness failed: %v", err))
			}

			resetTimer <- true
			// popup could be visible on other monitor, so we should hide it before revising its position, otherwise:
			// - it may cause a flicker because we're repositioning it, and
			// - we may also see brightness level of the previous monitor for some ms
			popupVisibleChan <- false
			popupPosChan <- m.getPosition()
			currentMonitor.name = m.getInstanceName()
			currentMonitor.brightness = b
			currentMonitor.monitor = &m
		} else {
			resetTimer <- true
		}
		popupVisibleChan <- true

		prevBrightness := currentMonitor.brightness
		switch command {
		case DecreaseBrightness:
			currentMonitor.brightness =
				clamp(0, 100,
					int(math.Floor(
						snapNumber(6.25)(float64(currentMonitor.brightness)-6.25))))
		case IncreaseBrightness:
			currentMonitor.brightness =
				clamp(0, 100,
					int(math.Floor(
						snapNumber(6.25)(float64(currentMonitor.brightness)+6.25))))
		}
		if currentMonitor.brightness != prevBrightness {
			go m.setBrightness(currentMonitor.brightness)
		}
	}
}

func clearCurrentMonitor(currentMonitor *MonitorInfo, resetTimer <-chan bool, popupVisibleChan chan<- bool) {
	t := time.AfterFunc(0, func() {})

	for {
		<-resetTimer
		t.Stop()
		t = time.AfterFunc(1*time.Second+200*time.Millisecond, func() {
			popupVisibleChan <- false
			currentMonitor.name = ""
			currentMonitor.brightness = 0
			currentMonitor.monitor = nil
		})
	}
}
