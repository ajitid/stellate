# Stellate

A Windows utility to control brightness of the display where your cursor is on using keyboard shortcuts.

Pause key to increase brightness  
Scroll Lock key to decrease brightness

_Stellate is in usable state but is not finished. That is why I haven't provided an .exe yet._

## How to build and use

1. Clone this repository
1. Install golang
1. Run `go build -ldflags "-H=windowsgui"` in the project root (this may take few mins to complete). A file called stellate.exe will be built.
1. Open Windows Explorer and copy this file. Then press <kbd>Win</kbd>+<kbd>R</kbd> key and type `shell:startup`. Right click and Paste Shortcut.
1. Double click on this shortcut. The program will run. 

You can close <img src="icon.ico" style="height: 1em; width: auto;" /> Stellate from system tray.  
The program will automatically run whenever you would restart Windows. If you don't want that you can delete that shortcut you pasted.
