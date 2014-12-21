package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpCodeCalculations(t *testing.T) {
	assert := assert.New(t)

	mem := make(memory, 2)
	mem[0] = 0xA2
	mem[1] = 0xF0
	op := opCode(0, mem)

	assert.Equal(op, 0xA2F0, "OpCode must be equal to 0xA2F0")
}
