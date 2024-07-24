// taken from https://claude.ai/chat/c1e10885-ed75-4106-96d3-7db09b64fd09
package main

import (
	"fmt"
	"log"
	"log/slog"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")

	dxva2                                   = syscall.NewLazyDLL("dxva2.dll")
	getNumberOfPhysicalMonitorsFromHMONITOR = dxva2.NewProc("GetNumberOfPhysicalMonitorsFromHMONITOR")
	getPhysicalMonitorsFromHMONITOR         = dxva2.NewProc("GetPhysicalMonitorsFromHMONITOR")
	getMonitorBrightness                    = dxva2.NewProc("GetMonitorBrightness")
	setMonitorBrightness                    = dxva2.NewProc("SetMonitorBrightness")
)

type (
	HMONITOR syscall.Handle

	POINT         struct{ X, Y int32 }
	RECT          struct{ Left, Top, Right, Bottom int32 }
	MONITORINFOEX struct {
		CbSize    uint32
		RcMonitor RECT
		RcWork    RECT
		DwFlags   uint32
		SzDevice  [32]uint16
	}
	PHYSICAL_MONITOR struct {
		hPhysicalMonitor             syscall.Handle
		szPhysicalMonitorDescription [128]uint16
	}
)

func cursorOnMonitor() (*HMONITOR, string, error) {
	var cursorPos POINT
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursorPos)))
	if ret == 0 {
		return nil, "", fmt.Errorf("GetCursorPos failed: %v", err)
	}

	// fmt.Printf("Cursor position: (%d, %d)\n", cursorPos.X, cursorPos.Y)

	var monitors []HMONITOR
	callback := syscall.NewCallback(func(hMonitor HMONITOR, hdcMonitor uintptr, lprcMonitor *RECT, dwData uintptr) uintptr {
		monitors = append(monitors, hMonitor)
		return 1 // continue enumeration
	})

	ret, _, err = procEnumDisplayMonitors.Call(0, 0, callback, 0)
	if ret == 0 {
		return nil, "", fmt.Errorf("EnumDisplayMonitors failed: %v", err)
	}

	for i, hMonitor := range monitors {
		var info MONITORINFOEX
		info.CbSize = uint32(unsafe.Sizeof(info))
		ret, _, err = procGetMonitorInfoW.Call(
			uintptr(hMonitor),
			uintptr(unsafe.Pointer(&info)),
		)
		if ret == 0 {
			slog.Warn(fmt.Sprintf("GetMonitorInfoW failed for monitor %d: %v\n", i, err))
			continue
		}

		deviceName := syscall.UTF16ToString(info.SzDevice[:])
		//fmt.Printf("Monitor %d: %s\n", i, deviceName)

		if cursorPos.X >= info.RcMonitor.Left && cursorPos.X < info.RcMonitor.Right &&
			cursorPos.Y >= info.RcMonitor.Top && cursorPos.Y < info.RcMonitor.Bottom {
			return &hMonitor, deviceName, nil
		}
	}

	return nil, "", fmt.Errorf("failed to find the monitor with the cursor")
}

type DDCMonitor HMONITOR // usually external displays

func (hMonitor DDCMonitor) setBrightness(value int) {
	var numPhysicalMonitors uint32
	ret, _, err := getNumberOfPhysicalMonitorsFromHMONITOR.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&numPhysicalMonitors)))
	if ret == 0 {
		log.Fatalf("GetNumberOfPhysicalMonitorsFromHMONITOR failed: %v\n", err)
	}

	physicalMonitors := make([]PHYSICAL_MONITOR, numPhysicalMonitors)
	ret, _, err = getPhysicalMonitorsFromHMONITOR.Call(
		uintptr(hMonitor),
		uintptr(numPhysicalMonitors),
		uintptr(unsafe.Pointer(&physicalMonitors[0])),
	)
	if ret == 0 {
		log.Fatalf("GetPhysicalMonitorsFromHMONITOR failed: %v\n", err)
	}

	for _, monitor := range physicalMonitors {
		ret, _, err = setMonitorBrightness.Call(uintptr(monitor.hPhysicalMonitor), uintptr(value))
		if ret == 0 {
			log.Fatalf("SetMonitorBrightness failed: %v\n", err)
		}
	}
}

func (hMonitor DDCMonitor) getBrightness() int {
	var numPhysicalMonitors uint32
	ret, _, err := getNumberOfPhysicalMonitorsFromHMONITOR.Call(uintptr(hMonitor), uintptr(unsafe.Pointer(&numPhysicalMonitors)))
	if ret == 0 {
		log.Fatalf("GetNumberOfPhysicalMonitorsFromHMONITOR failed: %v\n", err)
	}

	physicalMonitors := make([]PHYSICAL_MONITOR, numPhysicalMonitors)
	ret, _, err = getPhysicalMonitorsFromHMONITOR.Call(
		uintptr(hMonitor),
		uintptr(numPhysicalMonitors),
		uintptr(unsafe.Pointer(&physicalMonitors[0])),
	)
	if ret == 0 {
		log.Fatalf("GetPhysicalMonitorsFromHMONITOR failed: %v\n", err)
	}

	brightness := -1
	for _, monitor := range physicalMonitors {
		var minimumBrightness, currentBrightness, maximumBrightness uint32
		ret, _, err = getMonitorBrightness.Call(
			uintptr(monitor.hPhysicalMonitor),
			uintptr(unsafe.Pointer(&minimumBrightness)),
			uintptr(unsafe.Pointer(&currentBrightness)),
			uintptr(unsafe.Pointer(&maximumBrightness)),
		)
		if ret == 0 {
			log.Fatalf("GetMonitorBrightness failed: %v\n", err)
		}
		brightness = int(currentBrightness)
	}

	if brightness == -1 {
		log.Fatal("couldn't get brightness")
	}
	return brightness
}
