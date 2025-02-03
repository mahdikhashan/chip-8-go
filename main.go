package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

const (
	TITLE     = "CHIP8"
	WIDTH     = 640
	HEIGHT    = 480
	FRAMERATE = 60
)

const SCALE = 10
const TICKS_PER_FRAME = 10

var running = true

func main() {
	emu := initEmu()

	filename := os.Args[1]
	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	emu.load(dat)
	fmt.Print("data loaded into emulator")

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

		for range TICKS_PER_FRAME {
			emu.tick()
		}
		emu.tick_timer()
		draw_screen(emu, renderer)

		sdl.Delay(1000 / FRAMERATE)
	}
}

func draw_screen(emu *Emu, renderer *sdl.Renderer) {
	// Handle game logic (e.g., player movement, drawing objects, etc.)
	// Clear screen to black again for the next frame
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	var screen_buf = emu.get_display()
	renderer.SetDrawColor(255, 255, 255, 255)
	for i, pixel := range screen_buf {
		if pixel {
			x := uint32(uint(i) % SCREEN_WIDTH)
			y := uint32(uint(i) / SCREEN_WIDTH)
			rect := sdl.Rect{X: int32(x * SCALE), Y: int32(y * SCALE), W: SCALE, H: SCALE}
			err := renderer.FillRect(&rect)
			if err != nil {
				panic(err)
			}
		}
	}

	renderer.Present()
}

func key_to_button(key sdl.Keycode) {

}
