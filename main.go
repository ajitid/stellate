package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.design/x/hotkey"
)

const (
	WinWidth  = 180
	WinHeight = 40
)

var brightness = 0

func main() {
	/*
		Setting `FlagWindowHidden` before `InitWindow()` so that the window doesn't flashes (appears then quickly hides itself on start)
		Setting the rest in `SetWindowState()` as not every flag is configurable before window creation, see https://github.com/raysan5/raylib/issues/1367#issue-690893773
	*/
	rl.SetConfigFlags(rl.FlagWindowHidden | rl.FlagWindowTransparent | rl.FlagMsaa4xHint | rl.FlagWindowHighdpi)
	rl.InitWindow(WinWidth, WinHeight, "scintilla")
	defer rl.CloseWindow()
	/*
		I don't know if I _need_ HiDPI flag or not, see:
		https://github.com/raysan5/raylib/discussions/2999
		https://www.reddit.com/r/raylib/comments/o3k27c/macos_fix_high_dpi_blurry_window/
	*/
	rl.SetWindowState(rl.FlagWindowUndecorated | rl.FlagWindowTopmost)
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

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.NewColor(0, 0, 0, 0))

		// If rounded edge still look like chamfer, see https://www.reddit.com/r/raylib/comments/17obnui/how_to_anti_aliasing_for_shapes/ to improve
		// Right now I've used `rl.FlagMsaa4xHint` as config flag
		rl.DrawRectangleRounded(winSizedRect, 0.3, 8, rl.NewColor(1, 4, 9, 255))

		{
			// progress bg
			rl.DrawRectangle((WinWidth-progressWidth)/2+progressLeftPadding, (WinHeight-progressHeight)/2, progressWidth, progressHeight, rl.NewColor(27, 28, 32, 255))
		}

		{
			// progress fg
			rl.DrawRectangle((WinWidth-progressWidth)/2+progressLeftPadding, (WinHeight-progressHeight)/2, progressWidth*int32(brightness)/100, progressHeight, rl.NewColor(76, 194, 255, 255))
		}

		{
			var (
				x      int32   = WinWidth * 13 / 100
				y      int32   = WinHeight / 2
				radius float32 = 3
			)
			rl.DrawCircleLines(x, y, radius, rl.LightGray)
			drawLinesAroundCircle(rl.Vector2{X: float32(x), Y: float32(y)}, radius+4.3, 8, 2.1, rl.LightGray)
		}

		rl.EndDrawing()
	}
}

func registerHotkeys(brightnessCommandChan chan<- BrightnessCommand) {
	/*
		https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
		https://community.keyboard.io/t/what-are-the-codes-for-the-brightness-control-keys/4094
		https://www.reddit.com/r/Keychron/comments/1034z92/brightness_keys_mac_external_display/
	*/
	brightnessUpKey := hotkey.New([]hotkey.Modifier{}, hotkey.Key(0x13)) // win virtual keycode for pause
	err := brightnessUpKey.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return
	}
	defer brightnessUpKey.Unregister()

	brightnessDownKey := hotkey.New([]hotkey.Modifier{}, hotkey.Key(0x91)) // win virtual keycode for scroll lock
	err = brightnessDownKey.Register()
	if err != nil {
		log.Fatalf("hotkey: failed to register hotkey: %v", err)
		return
	}
	defer brightnessDownKey.Unregister()

	for {
		select {
		case <-brightnessUpKey.Keydown():
			brightnessCommandChan <- IncreaseBrightness
		case <-brightnessDownKey.Keydown():
			brightnessCommandChan <- DecreaseBrightness
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
