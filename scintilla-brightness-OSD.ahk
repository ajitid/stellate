#Requires AutoHotkey v2+
#SingleInstance Ignore

Cleanup(ExitReason, ExitCode) {
  Runwait "taskkill /im scintilla.exe", , "Hide"
}
OnExit(Cleanup)

; Commenting SetWorkingDir as it isn't needed if this .ahk is at same place where scintilla.exe is:
; SetWorkingDir "C:\Users\ajits\ghq\github.com\ajitid\scintilla"

; we'll kill any manually opened instances of scintilla or instances opened by other .ahk first
Runwait "taskkill /im scintilla.exe", , "Hide"
RunWait "scintilla.exe", , "Hide"