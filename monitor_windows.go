// taken from https://claude.ai/chat/b18f4a03-343f-479c-a0ce-7cb40c4404aa which I took from https://github.com/posthumz/DisplayDevices/blob/master/DisplayDevices.cs via this convo https://old.reddit.com/r/PowerShell/comments/19e7das/getting_display_id_for_a_display_device/
package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// DISPLAY_DEVICEW represents the DISPLAY_DEVICEW structure from Windows API
type DISPLAY_DEVICEW struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

func (dd *DISPLAY_DEVICEW) DeviceNameString() string {
	return syscall.UTF16ToString(dd.DeviceName[:])
}

func (dd *DISPLAY_DEVICEW) DeviceStringString() string {
	return syscall.UTF16ToString(dd.DeviceString[:])
}

func (dd *DISPLAY_DEVICEW) DeviceIDString() string {
	return syscall.UTF16ToString(dd.DeviceID[:])
}

func (dd *DISPLAY_DEVICEW) DeviceKeyString() string {
	return syscall.UTF16ToString(dd.DeviceKey[:])
}

/*
alternatively could've used a generic form of

// ToString converts a [32]uint16 to a string
func toString(buf [32]uint16) string {
	return syscall.UTF16ToString(buf[:])
}

// ToStringLong converts a [128]uint16 to a string
func toStringLong(buf [128]uint16) string {
	return syscall.UTF16ToString(buf[:])
}

usage:
deviceName := toString(devices[i].DeviceName)
*/

var (
	procEnumDisplayDevicesW = user32.NewProc("EnumDisplayDevicesW")
)

// EnumDisplayDevicesW wraps the Windows API function
func EnumDisplayDevicesW(device *uint16, devNum uint32, displayDevice *DISPLAY_DEVICEW, flags uint32) bool {
	ret, _, _ := procEnumDisplayDevicesW.Call(
		uintptr(unsafe.Pointer(device)),
		uintptr(devNum),
		uintptr(unsafe.Pointer(displayDevice)),
		uintptr(flags),
	)
	return ret != 0
}

// GetAll returns all display devices
func getDisplays(getInterfaceName uint32) []DISPLAY_DEVICEW {
	var devices []DISPLAY_DEVICEW

	for i := uint32(0); ; i++ {
		var dd DISPLAY_DEVICEW
		dd.Cb = uint32(unsafe.Sizeof(dd))

		if !EnumDisplayDevicesW(nil, i, &dd, getInterfaceName) {
			break
		}

		devices = append(devices, dd)

		for j := uint32(0); ; j++ {
			var monitor DISPLAY_DEVICEW
			monitor.Cb = uint32(unsafe.Sizeof(monitor))

			if !EnumDisplayDevicesW((*uint16)(unsafe.Pointer(&dd.DeviceName[0])), j, &monitor, getInterfaceName) {
				break
			}

			devices = append(devices, monitor)
		}
	}

	return devices
}

// FromID returns a display device by ID
func getDisplayFromID(id string, getInterfaceName uint32) *DISPLAY_DEVICEW {
	devices := getDisplays(getInterfaceName)
	for i := range devices {
		if strings.HasPrefix(devices[i].DeviceIDString(), fmt.Sprintf(`MONITOR\%s\`, id)) {
			return &devices[i]
		}
	}
	return nil
}

func getDisplayFromName(name string, getInterfaceName uint32) *DISPLAY_DEVICEW {
	devices := getDisplays(getInterfaceName)
	for i := range devices {
		if devices[i].DeviceNameString() == name {
			return &devices[i]
		}
	}
	return nil
}
