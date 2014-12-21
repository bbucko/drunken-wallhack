package chip8

import "log"

type memory []byte

type CPU struct {
	PC          uint16 `PC`
	I           uint16 `I`
	V           memory `V`
	draw        bool   `Draw`
	sound_timer byte   `sound_timer`
	delay_timer byte   `delay_timer`
	sp          byte   `sp`
	key         memory
	stack       memory
	mem         memory
	gfx         memory
}

func New() *CPU {
	c := &CPU{
		PC:          0,
		I:           0,
		draw:        false,
		sound_timer: 0,
		delay_timer: 0,
		V:           make(memory, 16),
		mem:         make(memory, 4096),
		gfx:         make(memory, 64*32),
		stack:       make(memory, 16),
		key:         make(memory, 16),
	}

	return c
}

func (c *CPU) Load(rom memory) error {
	copy(c.mem, rom)
	return nil
}

func (c *CPU) Step() error {
	memAddress := c.PC * 2
	opCode := opCode(memAddress, c.mem)
	log.Printf("%x", opCode)
	c.PC += 1
	return nil
}

func opCode(memAddress uint16, mem memory) uint16 {
	opCode := uint16(mem[memAddress])<<8 | uint16(mem[memAddress+1])
	log.Printf("%.16b %.8b %.8b", opCode, mem[memAddress], mem[memAddress+1])
	return opCode
}
