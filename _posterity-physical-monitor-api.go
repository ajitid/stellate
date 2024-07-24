package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	dxva2               = syscall.NewLazyDLL("dxva2.dll")
	enumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")

	getNumberOfPhysicalMonitorsFromHMONITOR = dxva2.NewProc("GetNumberOfPhysicalMonitorsFromHMONITOR")
	getPhysicalMonitorsFromHMONITOR         = dxva2.NewProc("GetPhysicalMonitorsFromHMONITOR")
	getMonitorBrightness                    = dxva2.NewProc("GetMonitorBrightness")
	setMonitorBrightness                    = dxva2.NewProc("SetMonitorBrightness")
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type PHYSICAL_MONITOR struct {
	hPhysicalMonitor             syscall.Handle
	szPhysicalMonitorDescription [128]uint16
}

func main() {
	err := enumDisplayMonitors.Find()
	if err != nil {
		fmt.Printf("Failed to find EnumDisplayMonitors: %v\n", err)
		return
	}

	err = enumDisplayMonitors.Call(0, 0, syscall.NewCallback(monitorEnumProc), 0)
	if err != nil && err != syscall.Errno(0) {
		fmt.Printf("EnumDisplayMonitors failed: %v\n", err)
	}
}

func monitorEnumProc(hMonitor syscall.Handle, hdcMonitor syscall.Handle, lprcMonitor *RECT, dwData uintptr) uintptr {
	var numPhysicalMonitors uint32
	ret, _, err := getNumberOfPhysicalMonitorsFromHMONITOR.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&numPhysicalMonitors)))
	if ret == 0 {
		fmt.Printf("GetNumberOfPhysicalMonitorsFromHMONITOR failed: %v\n", err)
		return 1
	}

	physicalMonitors := make([]PHYSICAL_MONITOR, numPhysicalMonitors)
	ret, _, err = getPhysicalMonitorsFromHMONITOR.Call(
		uintptr(hMonitor),
		uintptr(numPhysicalMonitors),
		uintptr(unsafe.Pointer(&physicalMonitors[0])),
	)
	if ret == 0 {
		fmt.Printf("GetPhysicalMonitorsFromHMONITOR failed: %v\n", err)
		return 1
	}

	for _, monitor := range physicalMonitors {
		var minimumBrightness, currentBrightness, maximumBrightness uint32
		ret, _, err = getMonitorBrightness.Call(
			uintptr(monitor.hPhysicalMonitor),
			uintptr(unsafe.Pointer(&minimumBrightness)),
			uintptr(unsafe.Pointer(&currentBrightness)),
			uintptr(unsafe.Pointer(&maximumBrightness)),
		)
		if ret == 0 {
			fmt.Printf("GetMonitorBrightness failed: %v\n", err)
			continue
		}

		fmt.Printf("Current brightness: %d\n", currentBrightness)

		// Increase brightness by 10%
		newBrightness := currentBrightness + 10
		if newBrightness > maximumBrightness {
			newBrightness = maximumBrightness
		}

		ret, _, err = setMonitorBrightness.Call(uintptr(monitor.hPhysicalMonitor), uintptr(newBrightness))
		if ret == 0 {
			fmt.Printf("SetMonitorBrightness failed: %v\n", err)
		} else {
			fmt.Printf("New brightness: %d\n", newBrightness)
		}
	}

	return 1
}
