**SDL**


## Installing

https://youtu.be/OXSMx45kayw?si=gGHnLd15Zin8-ov1

Steps:
1. Grab the latest release of mingw-w64 of https://github.com/niXman/mingw-builds-binaries and put the extracted folder into `C:\` so that it becomes `C:\mingw64` (I used x86_64-13.2.0-release-win32-seh-ucrt-rt_v11-rev0.7z)
1. Download `SDL2-devel-[version]-mingw.zip` from https://github.com/libsdl-org/SDL/releases. Extract it. Copy SDL's x86_64-w64-mingw32 folder contents into C:\mingw64\x86_64-w64-mingw32. The contents will be merged.
1. Edit the _system_ environment variables. Add C:\mingw64\x86_64-w64-mingw32\bin and C:\mingw64\bin to the path
1. Restart any open terminal (incl. VS Code)
1. In your project, run `go get github.com/veandco/go-sdl2/sdl` followed by `go build github.com/veandco/go-sdl2/sdl`. This will take few mins. (Prefer to run these commands in an actual Windows Terminal and not in VS Code as VS Code restores terminal sessions and you mistakenly run the command there.)
1. Then use this code:

```golang
package main

import "github.com/veandco/go-sdl2/sdl"

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
}
```
and run it using `go run .`. If it ends without printing anything, then it means it ran successfully. If you'd reopen VS Code you'd now see typings.  
VS Code will pick up typings eventually. I don't know what caused it but they became availble after few VS Code restarts.  
I noticed that using using Windows Terminal's Powershell to open code using `code .` always loads the typings.

## Attaching VS Code debugger

https://youtu.be/yxK_dwJ3Bbc?si=8uiDtvJacm5x_Y18&t=1000

## Getting release (statically linked) libs:

Download and install SDL2 runtime libraries from https://github.com/libsdl-org/SDL/releases. Extract and copy the .dll file into the project directory. After that, the program should become runnable.

Alternatively, static DLLs can be grabbed from here: https://github.com/mmozeiko/build-sdl2

---

**Raylib**

This may affect your window positioning logic https://github.com/raysan5/raylib/pull/3950
Also see https://stackoverflow.com/questions/64518580/borderless-window-covers-taskbar

## Using Raylib without CGo

https://github.com/gen2brain/raylib-go?tab=readme-ov-file#purego-without-cgo-ie-cgo_enabled0  
I put the DLL in my project root because that's where I usually put the `go build` output when testing. Also `go run .` picks up DLL from project's root too. Wherever you put it, just make sure to accompany your .exe with raylib.dll

## Other interesting libs

From https://github.com/gopxl/pixel/wiki/Creating-a-Window#run

> There's one ugly thing about graphics and operating systems. That one thing is that most operating systems require all graphics and windowing code to be executed from the main thread of our program. This is really cumbersome with Go. Go is a highly concurrent language with goroutines. Goroutines can freely jump from thread to thread, which makes the previous requirement seemingly impossible. Not all is lost, however. Go's runtime provides a convenient function runtime.LockOSThread, which locks current goroutine on it's current thread. PixelGL uses this functionality and provides you with a simpler interface to it.
>
> You don't have to deal with the main thread stuff at all with Pixel. You can run your game concurrently, however you want. You only need to allow Pixel to use the main thread.

- https://github.com/gopxl/mainthread for solving the above problem
- https://github.com/gopxl/beep for sound
- https://github.com/gopxl/glhf for OpenGL
- what is font atlas: https://github.com/gopxl/pixel/wiki/Typing-text-on-the-screen

## Transform origin

https://www.reddit.com/r/raylib/comments/15021nl/rectangle_not_at_the_center_of_the_screen_for/

## Font rendering

If you want to render text, use FilterBilinear with this trick https://www.raylib.com/examples/text/loader.html?name=text_font_filters

Also, use a higher font size https://github.com/raysan5/raylib/issues/2355#issuecomment-1050059197.
Update: ^ this higher font size trick didn't work for me. I rather specify the font size I need in LoadFontEx itself.
LoadFont() generates a font atlas with size of 32pixel height each glyph

Another thing to know is if you open a font file you'd see font displayed
in different sizes (some multiples): 12, 18, 24, 36, 48... Turns out the
fonts will be best displayed in these sizes. 
https://www.reddit.com/r/raylib/comments/1dqwldb/can_i_render_text_with_a_sdf_font_shader_to_a/lb02ld0/

See https://www.reddit.com/r/raylib/comments/xfrv7y/text_kerning/
https://www.reddit.com/r/typography/comments/vhglmd/how_to_get_proper_kerning_tracking_and_leading_in/
https://www.reddit.com/r/raylib/comments/12whhpw/is_it_possible_to_reduce_the_amount_of_space_in/

See other techniques: 
- https://www.raylib.com/examples/text/loader.html?name=text_font_sdf
- https://www.raylib.com/examples/text/loader.html?name=text_font_loading
	- https://www.reddit.com/r/raylib/comments/1b96jcb/how_to_use_bitmap_fonts/
- there are probably more techniques
- https://github.com/raysan5/raylib/issues/908

There's this but I don't know how helpful it is given I'm using go and I can get these values already https://www.reddit.com/r/raylib/comments/190taie/distributing_fonts_and_other_assets/

SDL font rendering will always give better quality though:
- https://www.reddit.com/r/raylib/comments/xfrv7y/text_kerning/
- https://gist.github.com/raysan5/17392498d40e2cb281f5d09c0a4bf798#file-formats-support
- https://github.com/raysan5/raylib/issues/2355#issuecomment-1050059197
- https://x.com/ludwigABAP/status/1817535771085074755

```go
// https://github.com/adrg/sysfont
fontFinder := sysfont.NewFinder(&sysfont.FinderOpts{
	Extensions: []string{".ttf"},
})
matchedFont := fontFinder.Match("Segoe UI")
if matchedFont == nil {
	log.Fatal("system font not found")
}
font := rl.LoadFontEx(matchedFont.Filename, 18, []rune(nil)) // Pass nil to load the default character set
defer rl.UnloadFont(font)

// Following code is not needed if you aren't using FilterTrilinear. On 2D drawing FilterTrilinear isn't noticeable because it looks like FilterBiLinear
// rl.GenTextureMipmaps(&font.Texture)
// rl.SetTextureFilter(font.Texture, rl.FilterTrilinear)
```

`GenTextureMipmaps` will give you blurry font if you'd use LoadFont() and then DrawTextEx with it on anything font size < 24.

Centering text:

```go
if loading {
	text := "Retrieving brightness..."
	var spacing float32 = 1
	textSize := rl.MeasureTextEx(font, text, fontSize, spacing)
	centeredPosition := rl.Vector2{
		X: float32(WinWidth)/2 - textSize.X/2,
		Y: float32(WinHeight)/2 - textSize.Y/2, // -1 to optically center it
	}
	rl.DrawTextEx(font, text, centeredPosition, fontSize, spacing, rl.LightGray)

	rl.EndDrawing()
	continue
}
```
