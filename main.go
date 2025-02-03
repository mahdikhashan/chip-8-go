package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	TITLE     = "CHIP8"
	WIDTH     = 640
	HEIGHT    = 480
	FRAMERATE = 60
)

var running = true

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(0, 0, 0, 255) // Black color
	renderer.Clear()

	renderer.Present()

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					}
				}
			}
		}

		// Handle game logic (e.g., player movement, drawing objects, etc.)
		// Clear screen to black again for the next frame
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		renderer.Present()

		sdl.Delay(1000 / FRAMERATE)
	}
}
