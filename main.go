package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.design/x/hotkey"
)

const (
	WinWidth  = 800
	WinHeight = 600
)

func main() {
	var brightnessCommandChan = make(chan BrightnessCommand, 1)
	go brightnessSetter(brightnessCommandChan)
	go registerHotkeys(brightnessCommandChan)

	rl.InitWindow(800, 450, "scintilla")
	rl.SetWindowState(rl.FlagWindowUndecorated)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)

		rl.EndDrawing()
	}
}

func registerHotkeys(brightnessCommandChan chan<- BrightnessCommand) {
	// https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
	// https://community.keyboard.io/t/what-are-the-codes-for-the-brightness-control-keys/4094
	// https://www.reddit.com/r/Keychron/comments/1034z92/brightness_keys_mac_external_display/
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
