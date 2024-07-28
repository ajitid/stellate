package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutex = kernel32.NewProc("CreateMutexW")
)

// from https://claude.ai/chat/38e56e68-e64a-4a1b-8272-7ac1a5e7ba82
func checkSingleInstance() bool {
	mutexName, err := syscall.UTF16PtrFromString("Global\\StellateBrightnessUtilityMutex")
	if err != nil {
		fmt.Println("Error creating mutex name:", err)
		return false
	}

	handle, _, lastErr := procCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(mutexName)),
	)

	if handle == 0 {
		fmt.Println("Error creating mutex:", lastErr)
		return false
	}

	// Check if the mutex already exists
	if lastErr == syscall.ERROR_ALREADY_EXISTS {
		syscall.CloseHandle(syscall.Handle(handle))
		fmt.Println("Another instance is already running")
		return false
	}

	// Keep the mutex handle open until the program exits
	return true
}
