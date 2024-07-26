This may affect your positioning logic https://github.com/raysan5/raylib/pull/3950

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