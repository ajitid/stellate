# SDL

## Installing

https://youtu.be/OXSMx45kayw?si=gGHnLd15Zin8-ov1
https://youtu.be/yxK_dwJ3Bbc?si=kiHyfIgx_D-85wu9

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

## Getting release (statically linked) libs:

Download and install SDL2 runtime libraries from https://github.com/libsdl-org/SDL/releases. Extract and copy the .dll file into the project directory. After that, the program should become runnable.

Alternatively, static DLLs can be grabbed from here: https://github.com/mmozeiko/build-sdl2

## Basics

https://www.youtube.com/watch?v=2VVFcs8jHRk&list=PLOXvU5Ov-cqpjd1_OnczdizY0I64OfH-T&index=2 
https://www.youtube.com/watch?v=OXSMx45kayw&list=PLDZujg-VgQlZUy1iCqBbe5faZLMkA3g2x&index=7
https://www.youtube.com/watch?v=2rs-LD9s7Js