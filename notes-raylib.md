This may affect your positioning logic https://github.com/raysan5/raylib/pull/3950

## Using Raylib without CGo

https://github.com/gen2brain/raylib-go?tab=readme-ov-file#purego-without-cgo-ie-cgo_enabled0  
I put the DLL in my project root because that's where I usually put the `go build` output when testing. Also `go run .` picks up DLL from project's root too. Wherever you put it, just make sure to accompany your .exe with raylib.dll