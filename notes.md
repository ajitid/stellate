## Embed .ico and manifest in Go program

https://github.com/akavel/rsrc
There's also https://gobyexample.com/embed-directive

## Display name to device ID mapping

[This](https://github.com/posthumz/DisplayDevices) from [here](https://www.reddit.com/r/PowerShell/comments/19e7das/getting_display_id_for_a_display_device/)

## Packages

Other packages to try:

```
"github.com/go-ole/go-ole"
"github.com/lxn/win"
"github.com/yusufpapurcu/wmi"
"github.com/winlabs/gowin32"
```

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

```powershell
Get-WmiObject -Query "SELECT DeviceID, Name FROM Win32_PnPEntity WHERE PNPClass = 'Monitor'"
```

## Softwares to control brightness

https://github.com/xanderfrangos/twinkle-tray
https://github.com/emoacht/Monitorian
https://www.nirsoft.net/utils/control_my_monitor.html
https://github.com/chrismah/ClickMonitorDDC7.2

## Getting a rnadom number

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