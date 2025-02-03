package main

import "math/bits"

const (
	RAM_SIZE      uint = 4096
	SCREEN_HEIGHT uint = 32
	SCREEN_WIDTH  uint = 64
	NUMS_REGS     uint = 16
	STACK_SIZE    uint = 16
	NUM_KEYS      uint = 16
	FONTSET_SIZE  uint = 80
)

var FONTSET = [FONTSET_SIZE]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Emu struct {
	pc     uint16
	ram    [RAM_SIZE]uint8
	screen [SCREEN_HEIGHT * SCREEN_WIDTH]bool
	v_reg  [NUMS_REGS]uint8
	i_reg  uint16
	sp     uint16
	stack  [STACK_SIZE]uint16
	keys   [NUM_KEYS]bool
	dt     uint8
	st     uint8
}

const START_ADDR uint16 = 0x200

func initEmu() Emu {
	e := Emu{
		pc: START_ADDR,
	}
	// i dono how this following like works
	copy(e.ram[:FONTSET_SIZE], FONTSET[:])
	return e
}

func (e *Emu) push(val uint16) {
	e.stack[e.sp] = val
	e.sp++
}

// TODO: handle underflow panic
func (e *Emu) pop() uint16 {
	e.sp--
	return e.stack[e.sp]
}

func (e *Emu) reset() {
	e.pc = START_ADDR
	e.ram = [RAM_SIZE]uint8{0}
	e.screen = [SCREEN_HEIGHT * SCREEN_WIDTH]bool{false}
	e.v_reg = [NUMS_REGS]uint8{0}
	e.i_reg = 0
	e.sp = 0
	e.stack = [STACK_SIZE]uint16{0}
	e.keys = [NUM_KEYS]bool{false}
	e.dt = 0
	e.st = 0
	// copy the fontset to the ram
	copy(e.ram[:FONTSET_SIZE], FONTSET[:])
}

func (e *Emu) tick() {
	var op = e.fetch()
	e.execute(op)
}

func (e *Emu) execute(op uint16) {
	var d1 = (op & 0xF000) >> 12
	var d2 = (op & 0x0F00) >> 8
	var d3 = (op & 0x00F0) >> 4
	var d4 = op & 0x000F

	var digits = [4]uint16{d1, d2, d3, d4}
	switch {
	// noop
	case digits == [4]uint16{0, 0, 0, 0}:
		return
	// clear screen
	case digits == [4]uint16{0, 0, 0xE, 0}:
		e.screen = [SCREEN_HEIGHT * SCREEN_WIDTH]bool{false}
	// return from subroutine
	case digits == [4]uint16{0, 0, 0xE, 0xE}:
		var ret_add = e.pop()
		e.pc = ret_add
	// jump
	case digits[0] == 1:
		var nnn = op & 0xFFF
		e.pc = nnn
	// call subroutine
	case digits[0] == 2:
		var nnn = op & 0xFFF
		e.push(e.pc)
		e.pc = nnn
	// skip
	case digits[0] == 3:
		var x = d2
		var nn = uint8(op & 0xFF)
		if e.v_reg[x] == nn {
			e.pc += 2
		}
	// skip if not equal VX !=Nn
	case digits[0] == 4:
		var x = d2
		var nn = uint8(op & 0xFF)
		if e.v_reg[x] != nn {
			e.pc += 2
		}
	case digits[0] == 5 && digits[3] == 0:
		var x = d2
		var y = d3
		if e.v_reg[x] == e.v_reg[y] {
			e.pc += 2
		}
	// 6XNN - VX = NN
	case digits[0] == 6:
		var x = d2
		var nn = uint8(op & 0xFF)
		e.v_reg[x] = nn
	case digits[0] == 7:
		var x = d2
		var nn = uint8(op & 0xFF)
		// wrt gpt, it can wrap automatically
		e.v_reg[x] = e.v_reg[x] + nn
	// 8XY0 - VX = VY
	case digits[0] == 8 && digits[3] == 0:
		var x = d2
		var y = d3
		e.v_reg[x] = e.v_reg[y]
	// TODO: 8XY1, 8XY2, 8XY3 - Bitwise operations
	case digits[0] == 8 && digits[3] == 1:
		var x = d2
		var y = d3
		// bitwise OR
		e.v_reg[x] |= e.v_reg[y]
	// 	seems to not correctly implemented
	case digits[0] == 8 && digits[3] == 4:
		var x = d2
		var y = d3
		newVx, carry := bits.Add(uint(e.v_reg[x]), uint(e.v_reg[y]), 0)
		var newVf uint
		if carry != 0 {
			newVf = 1
		} else {
			newVf = 0
		}
		e.v_reg[x] = uint8(newVx)
		e.v_reg[y] = uint8(newVf)
	default:
		panic("omg!")
	}
}

func (e *Emu) fetch() uint16 {
	var high_byte = uint16(e.ram[e.pc])
	var low_byte = uint16(e.ram[e.pc+1])
	// Big Endian
	var op = (high_byte << 8) | low_byte
	e.pc += 2
	return op
}

func (e *Emu) tick_timer() {
	if e.dt > 1 {
		e.dt -= 1
	}

	if e.st > 0 {
		if e.st == 1 {
			// TODO: implement beep
		}
		e.st -= 1
	}
}
