// from Claude https://claude.ai/chat/980e2e17-9922-4621-a63d-89e219679628
package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type MONITORINFOEX struct {
	win.MONITORINFO
	DeviceName [win.CCHDEVICENAME]uint16
}

type enumDisplayMonitorsCallback func(hMonitor win.HMONITOR, hdcMonitor win.HDC, lprcMonitor *win.RECT, dwData uintptr) uintptr

func enumDisplayMonitors(hdc win.HDC, lprcClip *win.RECT, lpfnEnum enumDisplayMonitorsCallback, dwData uintptr) bool {
	user32 := syscall.MustLoadDLL("user32.dll")
	defer user32.Release()

	enumDisplayMonitorsProc := user32.MustFindProc("EnumDisplayMonitors")

	ret, _, _ := enumDisplayMonitorsProc.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lprcClip)),
		syscall.NewCallback(lpfnEnum),
		dwData,
	)

	return ret != 0
}

func monitorEnumProc(hMonitor win.HMONITOR, hdcMonitor win.HDC, lprcMonitor *win.RECT, dwData uintptr) uintptr {
	var mi MONITORINFOEX
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	if win.GetMonitorInfo(hMonitor, (*win.MONITORINFO)(unsafe.Pointer(&mi))) {
		var cursorPos win.POINT
		win.GetCursorPos(&cursorPos)

		if cursorPos.X >= mi.RcMonitor.Left && cursorPos.X < mi.RcMonitor.Right &&
			cursorPos.Y >= mi.RcMonitor.Top && cursorPos.Y < mi.RcMonitor.Bottom {
			// Store the monitor handle
			*(*win.HMONITOR)(unsafe.Pointer(dwData)) = hMonitor
			return 0 // Stop enumeration
		}
	}

	return 1 // Continue enumeration
}

func getMonitorIdContainingCursor() (string, error) {
	var hMonitor win.HMONITOR

	enumDisplayMonitors(0, nil, monitorEnumProc, uintptr(unsafe.Pointer(&hMonitor)))

	if hMonitor != 0 {
		var mi MONITORINFOEX
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		if win.GetMonitorInfo(hMonitor, (*win.MONITORINFO)(unsafe.Pointer(&mi))) {
			// Convert device name to string
			monitorId := syscall.UTF16ToString(mi.DeviceName[:])
			return monitorId, nil
		} else {
			return "", fmt.Errorf("Failed to get monitor info.")
		}
	} else {
		return "", fmt.Errorf("Failed to find the monitor with the cursor.")
	}
}
