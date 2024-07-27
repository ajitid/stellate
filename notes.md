## Embed .ico and manifest in Go program

https://github.com/akavel/rsrc
There's also https://gobyexample.com/embed-directive

# Global hotkeys

https://pkg.go.dev/golang.design/x/hotkey
https://github.com/micmonay/keybd_event


## Display name to device ID mapping

[This](https://github.com/posthumz/DisplayDevices) from [here](https://www.reddit.com/r/PowerShell/comments/19e7das/getting_display_id_for_a_display_device/)

## Packages

Other packages to try:

```
"github.com/go-ole/go-ole"
"github.com/lxn/win"
"github.com/yusufpapurcu/wmi"
"github.com/winlabs/gowin32"
"github.com/hillu/go-ntdll"
"github.com/rodrigocfd/windigo"
"github.com/iamacarpet/go-win64api"
```

Code to reference:
- https://github.com/PoeBlu/hardentools/blob/master/autorun.go
- https://github.com/linexjlin/inputGPT
- https://github.com/st0le/winrec
- https://github.com/schollz/melrose

## Alternative way to get cursor position

Alternative way to get cursor in case Golang tells that `syscall.Syscall` is deprecated and recommends to use `syscall.SyscallN` instead:

```go
userDll := syscall.NewLazyDLL("user32.dll")
getWindowRectProc := userDll.NewProc("GetCursorPos")
type POINT struct {
	X, Y int32
}
var pt POINT
_, _, eno := syscall.SyscallN(getWindowRectProc.Addr(), uintptr(unsafe.Pointer(&pt)))
if eno != 0 {
	fmt.Println(eno)
}
fmt.Printf("[cursor.Pos] X:%d Y:%d\n", pt.X, pt.Y)
```

## Run command

```go
// get output
out, err := exec.Command("date").Output()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("The date is %s\n", out)

// just run + check for error
cmd := exec.Command("cmd.exe", "/C", "start", "/b", ".\\telltail-sync.ahk")
cmd.Dir = dir
if err := cmd.Run(); err != nil {
	log.Println("Error:", err)
}
```

## Get details of attached monitors

Option A:

```powershell
Get-WmiObject -Query "SELECT DeviceID, Name FROM Win32_PnPEntity WHERE PNPClass = 'Monitor'"
```

or option B:

```powershell
Get-WmiObject WmiMonitorID -Property InstanceName -Namespace root\wmi
# OR
Get-WmiObject WmiMonitorID -Namespace root\wmi | Where-Object { $_.InstanceName -like "DISPLAY\SOMESTRING\*" } | Select-Object -ExpandProperty InstanceName
# OR
Get-WmiObject -Query "SELECT InstanceName FROM WmiMonitorID" -Namespace root\wmi
```

Prefer option B as it doesn't capitalizes `14db058f` part of `DISPLAY\SHP1523\5&14db058f&2&UID512_0` and thus can directly be supplied to monitorian (after stripping _0 of course).

## Softwares to control brightness

https://github.com/xanderfrangos/twinkle-tray
https://github.com/emoacht/Monitorian
https://www.nirsoft.net/utils/control_my_monitor.html
https://github.com/chrismah/ClickMonitorDDC7.2

## Getting a random number

```go
// use "math/rand/v2" package

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

brightness := randRange(0, 101)
```

## Clean unused deps

```sh
go mod tidy
```

## Decoding WMI string

Just in case if I need it. See version 3 of 3 of the artifact here https://claude.ai/chat/03348c23-6712-4da4-bacb-7c682c935a14.
The struct in there suggest to use `string` type for bytes which is probably incorrect.

```golang
func decodeWMIString(s string) string {
	bytes := []byte(s)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 0 {
			return string(bytes[:i])
		}
	}
	return s
}
```

## Rest

https://github.com/nyaosorg/go-windows-su sudo/administrator
https://claude.ai/chat/46f64a64-2c6d-41ec-b77a-6f079bcd5662 get theme and theme change detection