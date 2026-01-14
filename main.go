package main

import (
	"fmt"
	"os"

	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	TITLE     = "CHIP8"
	WIDTH     = 64 * SCALE
	HEIGHT    = 32 * SCALE
	FRAMERATE = 60
)

const SCALE = 10
const TICKS_PER_FRAME = 10

var running = true

func init() {
	runtime.LockOSThread()
}

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
				k := key_to_button(e.Keysym.Sym)

				if k == 0xFF {
					continue
				}

				if e.Repeat > 0 {
					continue
				}

				if e.Type == sdl.KEYDOWN {
					// fmt.Println("in keydown")
					// k := key_to_button(e.Keysym.Sym)
					// fmt.Println(k)
					emu.keypress(uint(k), true)
				} else if e.Type == sdl.KEYUP {
					// fmt.Println("in keyup")
					// k := key_to_button(e.Keysym.Sym)
					// fmt.Println(k)
					emu.keypress(uint(k), false)
				}
			}
		}

		for range TICKS_PER_FRAME {
			emu.tick()
		}
		emu.tick_timer()
		draw_screen(emu, renderer)

		sdl.Delay(750 / FRAMERATE)
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

func key_to_button(key sdl.Keycode) uint16 {
	switch key {
	case sdl.K_1:
		return 0x1
	case sdl.K_2:
		return 0x2
	case sdl.K_3:
		return 0x3
	case sdl.K_4:
		return 0xC
	case sdl.K_q:
		return 0x4
	case sdl.K_w:
		return 0x5
	case sdl.K_e:
		return 0x6
	case sdl.K_r:
		return 0xD
	case sdl.K_a:
		return 0x7
	case sdl.K_s:
		return 0x8
	case sdl.K_d:
		return 0x9
	case sdl.K_f:
		return 0xE
	case sdl.K_z:
		return 0xA
	case sdl.K_x:
		return 0x0
	case sdl.K_c:
		return 0xB
	case sdl.K_v:
		return 0xF
	default:
		return 0xFF
	}
}
