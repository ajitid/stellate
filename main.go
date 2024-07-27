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
	startSystray, stopSystray := systray.RunWithExternalLoop(func() {
		iconBytes, err := os.ReadFile("icon.ico")
		if err != nil {
			log.Fatal(err)
		}
		systray.SetIcon(iconBytes)
		systray.SetTitle("Stellate")
	}, func() {})
	startSystray()
	go setupSystray()
	defer stopSystray()

	/*
		Setting `FlagWindowHidden` before `InitWindow()` so that the window doesn't flashes (appears then quickly hides itself on start)
		Setting the rest in `SetWindowState()` as not every flag is configurable before window creation, see https://github.com/raysan5/raylib/issues/1367#issue-690893773
	*/
	rl.SetConfigFlags(rl.FlagWindowHidden | rl.FlagWindowTransparent | rl.FlagMsaa4xHint | rl.FlagWindowHighdpi)
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

	rl.SetTargetFPS(60)

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
		progressWidth       int32 = (WinWidth * 0.7)
		progressLeftPadding int32 = 15
		progressHeight      int32 = 4
	)

	for !quit {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Blank)

		// If rounded edge still look like chamfer, see https://www.reddit.com/r/raylib/comments/17obnui/how_to_anti_aliasing_for_shapes/ to improve
		// Right now I've used `rl.FlagMsaa4xHint` as config flag
		rl.DrawRectangleRounded(winSizedRect, 0.3, 8, rl.NewColor(1, 4, 9, 255))

		// progress bg
		rl.DrawRectangle((WinWidth-progressWidth)/2+progressLeftPadding, (WinHeight-progressHeight)/2, progressWidth, progressHeight, rl.NewColor(27, 28, 32, 255))
		// progress fg
		rl.DrawRectangle((WinWidth-progressWidth)/2+progressLeftPadding, (WinHeight-progressHeight)/2, progressWidth*int32(currentMonitor.brightness)/100, progressHeight, rl.NewColor(76, 194, 255, 255))

		{
			var (
				x      int32   = WinWidth * 13 / 100
				y      int32   = WinHeight / 2
				radius float32 = 3
			)
			rl.DrawCircleLines(x, y, radius, rl.LightGray)
			drawLinesAroundCircle(rl.Vector2{X: float32(x), Y: float32(y)}, radius+4.3, 8, mapRange(float32(currentMonitor.brightness), 0, 100, 1, 3), rl.LightGray)
		}

		/*
			If you want to render text, use FilterBilinear with this trick
			https://github.com/raysan5/raylib/issues/2355#issuecomment-1050059197.

			Another thing to know is if you open a font file you'd see font displayed
			in different sizes (some multiples): 12, 18, 24, 36, 48... Turns out the
			fonts will be best displayed in these sizes. So when loading a font with
			LoadFontEx(), choose a higher size from this multiple and also when
			drawing the text, prefer to stick to a size that belongs to this multiple.
			https://www.reddit.com/r/raylib/comments/1dqwldb/can_i_render_text_with_a_sdf_font_shader_to_a/lb02ld0/

			SDL font rendering will always give better quality though:
			- https://www.reddit.com/r/raylib/comments/xfrv7y/text_kerning/
			- https://gist.github.com/raysan5/17392498d40e2cb281f5d09c0a4bf798#file-formats-support

			Finally, find locally installed system font (like Segoe UI) with this
			https://github.com/adrg/sysfont
		*/

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
