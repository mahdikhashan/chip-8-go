package main

const (
	RAM_SIZE      uint = 4096
	SCREEN_HEIGHT uint = 32
	SCREEN_WIDTH  uint = 64
	NUMS_REGS     uint = 16
	STACK_SIZE    uint = 16
)

type Emu struct {
	pc     uint16
	ram    [RAM_SIZE]uint8
	screen [SCREEN_HEIGHT * SCREEN_WIDTH]bool
	v_reg  [NUMS_REGS]uint8
	i_reg  uint16
	sp     uint16
	stack  [STACK_SIZE]uint16
}
