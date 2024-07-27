package main

import (
	"log"

	"golang.org/x/sys/windows"
)

// Windows API constants and types
const (
	GWL_EXSTYLE      = -20
	WS_EX_TOOLWINDOW = 0x00000080
)

var (
	// user32               = windows.NewLazySystemDLL("user32.dll")
	procGetWindowLongPtr = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
)

func getWindowLongPtr(hwnd windows.HWND, index int32) uintptr {
	ret, _, err := procGetWindowLongPtr.Call(uintptr(hwnd), uintptr(index))
	if ret == 0 {
		log.Fatalf("getWindowLongPtr failed: %v", err)
	}
	return ret
}

func setWindowLongPtr(hwnd windows.HWND, index int32, value uintptr) uintptr {
	ret, _, err := procSetWindowLongPtr.Call(uintptr(hwnd), uintptr(index), value)
	if ret == 0 {
		log.Fatalf("procSetWindowLongPtr failed: %v", err)
	}
	return ret
}
