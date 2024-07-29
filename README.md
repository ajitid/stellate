# Stellate

A Windows utility to control brightness of the display where your cursor is on using keyboard shortcuts.

## How to Download and Use

1. Download the [latest release](https://github.com/ajitid/stellate/releases/latest). It would be named `stellate-<version>.zip`. Extract it.
1. Inside the extracted folder there would be `stellate.exe`. Copy it, and then press <kbd>Win</kbd>+<kbd>R</kbd> keys and type `shell:startup`. Upon hitting enter a folder will open up. Right click and choose Paste Shortcut.
1. Double-click on this shortcut. The program will run. 


Use Pause key to increase and Scroll Lock key to decrease brightness.

You can close <img src="icon.ico" style="height: 1em; width: auto;" /> Stellate from system tray.  
The program will automatically run whenever you would restart Windows. If you don't want that you can delete that shortcut you pasted.

## How to Build and Use

1. Clone this repository
1. Install golang
1. Run `go build -ldflags "-H=windowsgui"` in the project root (this may take few mins to complete). A file called stellate.exe will be built.
1. You can double-click on it to run it and close it from the system tray