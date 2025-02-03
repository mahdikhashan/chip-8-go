package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InitsEmuCorrectly(t *testing.T) {
	assert.Equal(t, "emu", initEmu(), `should be "FizzBuzz".`)
}
