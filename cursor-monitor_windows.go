package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
)

type (
	POINT         struct{ X, Y int32 }
	RECT          struct{ Left, Top, Right, Bottom int32 }
	MONITORINFOEX struct {
		CbSize    uint32
		RcMonitor RECT
		RcWork    RECT
		DwFlags   uint32
		SzDevice  [32]uint16
	}
	HMONITOR uintptr
)

func cursorOnMonitor() {
	var cursorPos POINT
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&cursorPos)))
	if ret == 0 {
		fmt.Printf("GetCursorPos failed: %v\n", err)
		return
	}

	fmt.Printf("Cursor position: (%d, %d)\n", cursorPos.X, cursorPos.Y)

	var monitors []HMONITOR
	callback := syscall.NewCallback(func(hMonitor HMONITOR, hdcMonitor uintptr, lprcMonitor *RECT, dwData uintptr) uintptr {
		monitors = append(monitors, hMonitor)
		return 1 // continue enumeration
	})

	ret, _, err = procEnumDisplayMonitors.Call(0, 0, callback, 0)
	if ret == 0 {
		fmt.Printf("EnumDisplayMonitors failed: %v\n", err)
		return
	}

	for i, hMonitor := range monitors {
		var info MONITORINFOEX
		info.CbSize = uint32(unsafe.Sizeof(info))
		ret, _, err = procGetMonitorInfoW.Call(
			uintptr(hMonitor),
			uintptr(unsafe.Pointer(&info)),
		)
		if ret == 0 {
			fmt.Printf("GetMonitorInfoW failed for monitor %d: %v\n", i, err)
			continue
		}

		deviceName := syscall.UTF16ToString(info.SzDevice[:])
		fmt.Printf("Monitor %d: %s\n", i, deviceName)

		if cursorPos.X >= info.RcMonitor.Left && cursorPos.X < info.RcMonitor.Right &&
			cursorPos.Y >= info.RcMonitor.Top && cursorPos.Y < info.RcMonitor.Bottom {
			fmt.Printf("Mouse cursor is on monitor: %s\n", deviceName)
			break
		}
	}
}
