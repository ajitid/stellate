#Requires AutoHotkey v2+
#SingleInstance Ignore

Cleanup(ExitReason, ExitCode) {
  Runwait "taskkill /im stellate.exe", , "Hide"
}
OnExit(Cleanup)

; Commenting SetWorkingDir as it isn't needed if this .ahk is at same place where stellate.exe is:
; SetWorkingDir "C:\Users\ajits\ghq\github.com\ajitid\stellate"

; we'll kill any manually opened instances of stellate or instances opened by other .ahk first
Runwait "taskkill /im stellate.exe", , "Hide"
RunWait "stellate.exe", , "Hide"