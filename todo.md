- ✅ Store in systray https://github.com/fyne-io/systray + https://github.com/linexjlin/inputGPT/blob/main/go.mod#L10 
- ✅ Add manifest and ico https://github.com/akavel/rsrc
- use floral pattern for logo (see screenshot taken in photos app in phone)
- use github CI to auto build
- Pause key (to increase brightness) affects play/pause in MPV and thus in Plex.
- see syncthing docs to see all the possible options to start the app on startup
  - alternatively, systray mentions a method with go flags that doesn't popup cmd prompt (this is a go thing and so it works without using systray as well) https://github.com/fyne-io/systray?tab=readme-ov-file#windows
  - obviously you'd want to keep only one instance running. https://claude.ai/chat/38e56e68-e64a-4a1b-8272-7ac1a5e7ba82 with `taskkill /im stellate.exe` (preferably full path to stellate) may work.
  - for now, autohotkey method is ok 
  - or better:
    - embed stuff using akavel/rsrc, use flags mentioned in systray to command prompt doesn't open upon start, and finally use systray with an option to quit
    - add 2 bat scripts. One for install another for uninstall. Use with win+r>"shell:startup" to find Startup location (or see what i did in telltail). Use the name stellate-osd.exe so it is easier for install/uninstall script to taskkill. (yep taskkill stellate-osd before installing as well)
    - use https://claude.ai/chat/38e56e68-e64a-4a1b-8272-7ac1a5e7ba82 to ensure only 1 instance is running
- convert all float64 to float32 and avoid unnecessary casting. Use generics if needed
- filenames don't make sense. Rename them and group stuff properly.
- because it uses accelerated renderer, monitor disconnects causes issues (hybrid/optmised mode). Rounded borders go away for one. And switching from Standard to Optimised kills the app. Using SDL2 software renderer may help here. Otherwise running on dGPU may kill the battery. See https://www.reddit.com/r/raylib/comments/191o1xz/transparent_overlay_window_help/
  - sdl2 accelerated renderer may solve the rounded corners problem as it using its own windowing system and not GLFW