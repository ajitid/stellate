package main

import (
	"log"
	"os"

	"fyne.io/systray"
	rl "github.com/gen2brain/raylib-go/raylib"
	hook "github.com/robotn/gohook"
	"golang.org/x/sys/windows"
)

const (
	WinWidth  = 180
	WinHeight = 40
)

var quit = false // probably ok that it is not protected by a mutex

func main() {
	if !checkSingleInstance() {
		os.Exit(1)
	}

	startSystray, stopSystray := systray.RunWithExternalLoop(func() {
		iconBytes, err := os.ReadFile("icon.ico")
		if err != nil {
			log.Fatal(err)
		}
		systray.SetIcon(iconBytes)
		systray.SetTooltip("Stellate")
		systray.SetTitle("Stellate")
	}, func() {})
	startSystray()
	go setupSystray()
	defer stopSystray()

	/*
		Setting `FlagWindowHidden` before `InitWindow()` so that the window doesn't flashes (appears then quickly hides itself on start)
		Setting the rest in `SetWindowState()` as not every flag is configurable before window creation, see https://github.com/raysan5/raylib/issues/1367#issue-690893773
	*/
	rl.SetConfigFlags(rl.FlagWindowHidden | rl.FlagWindowTransparent | rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(WinWidth, WinHeight, "stellate")
	defer rl.CloseWindow()
	/*
		I don't know if I _need_ HiDPI flag or not, see:
		https://github.com/raysan5/raylib/discussions/2999
		https://www.reddit.com/r/raylib/comments/o3k27c/macos_fix_high_dpi_blurry_window/
	*/
	rl.SetWindowState(rl.FlagWindowUndecorated | rl.FlagWindowTopmost | rl.FlagWindowUnfocused)

	// hide window from showing up in the taskbar whenever `rl.FlagWindowHidden` flag is cleared
	hwnd := rl.GetWindowHandle()
	windowLong := getWindowLongPtr(windows.HWND(hwnd), GWL_EXSTYLE)
	setWindowLongPtr(windows.HWND(hwnd), GWL_EXSTYLE, windowLong|WS_EX_TOOLWINDOW)

	{
		brightnessCommandChan := make(chan BrightnessCommand, 1)
		popupVisibleChan := make(chan bool)
		popupPosChan := make(chan rl.Vector2)

		go brightnessSetter(brightnessCommandChan, popupVisibleChan, popupPosChan)
		go registerHotkeys(brightnessCommandChan)
		go handlePopupVisibility(popupVisibleChan, popupPosChan)
	}

	winSizedRect := rl.Rectangle{
		Width:  WinWidth,
		Height: WinHeight,
	}

	var (
		progressWidth  float32 = (WinWidth * 0.7)
		progressHeight float32 = 4
		progressX      float32 = (WinWidth-progressWidth)/2 + 15
		progressY      float32 = (WinHeight - progressHeight) / 2
	)
	progressBgRect := rl.Rectangle{
		X:      progressX,
		Y:      progressY,
		Width:  progressWidth,
		Height: progressHeight,
	}
	progressFgRect := rl.Rectangle{
		X:      progressX,
		Y:      progressY,
		Width:  0,
		Height: progressHeight,
	}

	for !quit {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Blank)

		// If rounded edge still look like chamfer, see https://www.reddit.com/r/raylib/comments/17obnui/how_to_anti_aliasing_for_shapes/ to improve
		// Right now I've used `rl.FlagMsaa4xHint` as config flag
		rl.DrawRectangleRounded(winSizedRect, 0.3, 8, rl.NewColor(1, 4, 9, 255))

		// progress bg
		rl.DrawRectangleRec(progressBgRect, rl.NewColor(27, 28, 32, 255))
		// progress fg
		progressFgRect.Width = progressWidth * float32(currentMonitor.brightness) / 100
		rl.DrawRectangleRec(progressFgRect, rl.NewColor(57, 166, 222, 255))

		{
			center := rl.Vector2{
				X: WinWidth * 13 / 100,
				Y: WinHeight / 2,
			}
			var radius float32 = 3
			rl.DrawCircleV(center, radius, rl.LightGray)
			lineLength := mapRange(float32(currentMonitor.brightness), 0, 100, 1, 3)
			drawLinesAroundCircle(center, radius+4.3, 8, lineLength, rl.LightGray)
		}

		rl.EndDrawing()
	}
}

func setupSystray() {
	name := systray.AddMenuItem("Stellate", "")
	name.Disable()
	systray.AddSeparator()
	exit := systray.AddMenuItem("Exit", "")

	<-exit.ClickedCh
	quit = true
}

func registerHotkeys(brightnessCommandChan chan<- BrightnessCommand) {
	evChan := hook.Start()
	defer hook.End()

	/*
		chosen scroll lock & pause because they are repurposed in macOS for brightness as well:
		https://community.keyboard.io/t/what-are-the-codes-for-the-brightness-control-keys/4094
		https://www.reddit.com/r/Keychron/comments/1034z92/brightness_keys_mac_external_display/
	*/
	for ev := range evChan {
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown {
			// fmt.Printf("Key event - Raw code: %v, Keychar: %c\n", ev.Rawcode, ev.Keychar)

			// The rawcode for pause and scroll lock is possibily limited to Windows. Other OSs may emit some other rawcode
			switch ev.Rawcode {
			case 19: // Pause key
				brightnessCommandChan <- IncreaseBrightness
			case 145: // Scroll Lock
				brightnessCommandChan <- DecreaseBrightness
			}
		}
	}
}

func handlePopupVisibility(popupVisibleChan <-chan bool, popupPosChan <-chan rl.Vector2) {
	for {
		select {
		case visible := <-popupVisibleChan:
			if visible {
				rl.ClearWindowState(rl.FlagWindowHidden)
			} else {
				rl.SetWindowState(rl.FlagWindowHidden)
			}
		case pos := <-popupPosChan:
			windowMonitor := -1
			for i := range rl.GetMonitorCount() {
				mPos := rl.GetMonitorPosition(i)
				if mPos.X == pos.X && mPos.Y == pos.Y {
					windowMonitor = i
				}
			}
			if windowMonitor == -1 {
				log.Fatal("monitor to set window to is not found.")
			}

			rl.SetWindowMonitor(windowMonitor)
			/*
				Windows positions the window at center, so for X we'll pass as is.
				For Y, pos.Y gives top position of the monitor wrt the overall monitor setup. So `pos.Y + GetMonitorHeight()` will actually give us the the bottom value of the monitor.
			*/
			rl.SetWindowPosition(int(rl.GetWindowPosition().X), int(pos.Y)+rl.GetMonitorHeight(windowMonitor)-(WinHeight+60))
		}
	}
}
