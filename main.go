package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func handleKeypress(ev *sdl.KeyboardEvent, brightnessCommandChan chan<- BrightnessCommand) {
	// for keyboard modifiers, use ev.Keysym.Mod (preferred) or sdl.GetModState()
	if ev.Type == sdl.KEYDOWN {
		if ev.Repeat == 0 {
			if ev.Keysym.Sym == sdl.K_LEFT {
				brightnessCommandChan <- DecreaseBrightness
			} else if ev.Keysym.Sym == sdl.K_RIGHT {
				brightnessCommandChan <- IncreaseBrightness
			}
		} else {

		}
	} else if ev.Type == sdl.KEYUP {

	}
}

func update() {
}

func main() {
	var brightnessCommandChan = make(chan BrightnessCommand, 1)
	go brightnessSetter(brightnessCommandChan)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	win, err := sdl.CreateWindow("testing sdl2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Destroy()

	rnr, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}
	defer rnr.Destroy()

	quit := false
	for !quit {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch ev := ev.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				handleKeypress(ev, brightnessCommandChan)
			}
		}

		rnr.SetDrawColor(128, 48, 122, 255)
		rnr.Clear()

		update()

		rnr.Present()

		// sdl.Delay(16) // Cap at roughly 60 FPS (claude gave this, there would be a better way to get 16.67 or something dynamically)
	}
}
