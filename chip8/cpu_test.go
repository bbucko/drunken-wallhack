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
	op := decodeOpCode(0, mem)

	assert.Equal(0xA2F0, op, "OpCode must be equal to 0xA2F0")
}

func TestCallSubroutine(t *testing.T) {
	assert := assert.New(t)
	c := New(memory([]byte{}))
	c.sp = 0
	c.PC = 0x200
	c.callSubroutine(0x2111)
	assert.Equal(1, c.sp, "StackPointer must be set to 1")
	assert.Equal(0x200, c.stack[0], "Stack at address 0 must be equal to previous PC (0x200)")
	assert.Equal(0x111, c.PC, "PC must be equal to 0x111")
}

func TestLoadToRegister(t *testing.T) {
	assert := assert.New(t)
	c := New(memory([]byte{}))
	c.loadToRegister(0x6011)
	assert.Equal(0x11, c.V[0], "Index must be set to 0x111")
	assert.Equal(0x202, c.PC, "Program counter must be increased by two")
}

func TestLoadToI(t *testing.T) {
	assert := assert.New(t)
	c := New(memory([]byte{}))
	c.loadToI(0xA111)
	assert.Equal(0x111, c.I, "I must be set to 0x111")
	assert.Equal(0x202, c.PC, "Program counter must be increased by two")
}

func TestReturnFromSubroutine(t *testing.T) {
	assert := assert.New(t)
	c := New(memory([]byte{}))
	c.PC = 0x206
	c.sp = 1
	c.stack[0] = 0x204
	c.returnFromSubroutine()
	assert.Equal(0x204, c.PC)
	assert.Equal(0, c.sp)
	assert.Equal(0x204, c.stack[0])
}

func prepareMemory(mem memory) (c *CPU) {
	c = New(mem)
	return
}
