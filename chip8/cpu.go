package chip8

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type memory []byte
type instruction uint16
type address uint16

type CPU struct {
	PC          address `PC`
	I           address `I`
	V           memory  `V`
	draw        bool    `Draw`
	sound_timer byte    `sound_timer`
	delay_timer byte    `delay_timer`
	sp          byte    `sp`
	stack       []address
	key         memory
	mem         memory
	gfx         memory
}

var BEGIN_LOAD = address(0x200)

var FONT memory = memory([]byte{
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
})

func New(rom memory) *CPU {
	c := &CPU{
		PC:          BEGIN_LOAD,
		I:           0,
		draw:        false,
		sound_timer: 0,
		delay_timer: 0,
		stack:       make([]address, 16),
		V:           make(memory, 16),
		mem:         make(memory, 4096),
		gfx:         make(memory, 64*32),
		key:         make(memory, 16),
	}
	copy(c.mem[c.PC:], rom)
	copy(c.mem[0x50:], FONT)
	rand.Seed(time.Now().UnixNano())

	return c
}

func (c *CPU) Step() error {
	opCode := decodeOpCode(c.PC, c.mem)

	instruction := (opCode & 0xF000) >> 12
	if instruction == 0x0 {
		//TODO Refactor
		if opCode == 0x00E0 {
			c.PC = c.PC + 2

			log.Println("%x CLS", c.PC)
		} else if opCode == 0x00EE {
			c.returnFromSubroutine()
		} else {
			c.PC = c.PC + 2

			log.Printf("%x SYS %x", c.PC, opCode&0xFF0F)
		}
	} else if instruction == 0x1 {
		//Jump to address nnn
		c.jumpToAddress(opCode)
	} else if instruction == 0x2 {
		//Execute subroutine starting at address NNN
		c.callSubroutine(opCode)
	} else if instruction == 0x3 {
		//SE Vx, byte
		//TODO Refactor
		registry := int(opCode & 0x0F00 >> 8)
		value := byte(opCode & 0x00FF)
		if c.V[registry] == value {
			c.PC = c.PC + 2
		}
		c.PC = c.PC + 2
		log.Printf("%x SE V%d, %.2x", c.PC, registry, value)
	} else if instruction == 0x4 {
		//SNE Vx, byte
		//TODO Refactor
		registry := int(opCode & 0x0F00 >> 8)
		value := byte(opCode & 0x00FF)
		if c.V[registry] != value {
			c.PC = c.PC + 2
		}
		c.PC = c.PC + 2
		log.Printf("%x SNE V%d, %.2x", c.PC, registry, value)
	} else if instruction == 0x6 {
		//LD Vx, byte
		c.loadToRegister(opCode)
	} else if instruction == 0x7 {
		//ADD Vx, byte
		//TODO Refactor
		x := uint8(opCode & 0x0F00 >> 8)
		value := byte(opCode & 0x00FF)
		c.V[x] = c.V[x] + value
		c.PC = c.PC + 2
		log.Printf("%x ADD V%d, %.2x", c.PC, x, value)
	} else if instruction == 0x8 {
		//TODO Refactor
		if opCode&0xF00F == 0x8000 {
			x := opCode & 0x0F00 >> 8
			y := opCode & 0x00F0 >> 4
			log.Printf("%x LD V%d, V%d", c.PC, x, y)
		} else if opCode&0xF00F == 0x8001 {
			//TODO
		} else if opCode&0xF00F == 0x8002 {
			//TODO
		} else if opCode&0xF00F == 0x8003 {
			//TODO
		} else if opCode&0xF00F == 0x8004 {
			c.addRegisters(opCode)
		} else if opCode&0xF00F == 0x8005 {
			//TODO
		} else if opCode&0xF00F == 0x8006 {
			//TODO
		} else if opCode&0xF00F == 0x8007 {
			//TODO
		} else if opCode&0xF00F == 0x800E {
			//TODO
		} else {
			log.Fatalf("Unknown 0x800 instruction %x", opCode)
		}
		c.PC = c.PC + 2
	} else if instruction == 0xA {
		//LD I, addr
		c.loadToI(opCode)
	} else if instruction == 0xC {
		//Sets VX to a random number and NN
		//TODO Refactor
		x := uint8(opCode & 0x0F00)
		random := uint8(rand.Intn(255))
		c.V[x] = random & uint8(opCode&0x00FF)
		log.Printf("%x RND V%d, %d", c.PC, x, c.V[x])
		c.PC = c.PC + 2
	} else if instruction == 0xD {
		//DRW Vx, Vy, nibble
		//TODO Refactor
		x := opCode & 0x0F00 >> 8
		y := opCode & 0x00F0 >> 4
		n := opCode & 0x000F
		c.PC = c.PC + 2
		log.Printf("%x DRW V%d, V%d, %d", c.PC, x, y, n)
	} else if instruction == 0xF {
		//TODO Refactor
		registry := uint8(opCode & 0x0F00 >> 8)
		operation := opCode & 0xF0FF
		if operation == 0xF007 {
			log.Printf("%x Set V[%d] = delay timer value", c.PC, registry)
			c.V[registry] = c.delay_timer
		} else if operation == 0xF015 {
			log.Printf("%x Set delay timer value = V[%d]", c.PC, registry)
			c.delay_timer = c.V[registry]
		} else {
			log.Printf("%x F %x %x", c.PC, opCode, operation)
		}
		c.PC = c.PC + 2
	} else {
		log.Fatalf("%x Unknown %x %x", c.PC, opCode, opCode>>12)
	}

	if c.sound_timer > 0 {
		if c.sound_timer == 1 {
			fmt.Println("BEEP")
		}
		c.sound_timer--
	}

	if c.delay_timer > 0 {
		c.delay_timer--
	}
	return nil
}

func decodeOpCode(memAddress address, mem memory) instruction {
	opCode := instruction(mem[uint16(memAddress)])<<8 | instruction(mem[uint16(memAddress+1)])
	return instruction(opCode)
}

func (c *CPU) callSubroutine(opCode instruction) {
	//Call subroutine at nnn
	addr := address(opCode & 0x0FFF)
	c.stack[c.sp] = c.PC
	c.sp++
	c.PC = addr

	log.Printf("%x Call subroutine at %.3x %x", c.PC, addr, c.PC)
}

func (c *CPU) loadToRegister(opCode instruction) {
	//LD Vx, byte
	registerId := uint8(opCode & 0x0F00 >> 8)
	byte := byte(opCode & 0x00FF)
	c.V[registerId] = byte
	c.PC = c.PC + 2

	log.Printf("%x LD V%d, %.2x", c.PC, registerId, byte)
}

func (c *CPU) loadToI(opCode instruction) {
	addr := address(opCode & 0x0FFF)
	c.I = addr
	c.PC = c.PC + 2

	log.Printf("%x LD I, %.3x", c.PC, addr)
}

func (c *CPU) jumpToAddress(opCode instruction) {
	c.PC = address(opCode & 0x0FFF)
	log.Printf("%x JP %.3x", c.PC, c.PC)
}

func (c *CPU) returnFromSubroutine() {
	c.sp--
	c.PC = c.stack[c.sp]

	log.Printf("%x RET %.3x %d", c.PC, c.PC, c.sp)
}

func (c *CPU) addRegisters(opCode instruction) {
	//ADD Vx, Vy
	x := uint8(opCode & 0x0F00 >> 8)
	y := uint8(opCode & 0x00F0 >> 4)
	z := uint16(c.V[x]) + uint16(c.V[y])
	if z > 255 {
		c.V[0xf] = 1
	}
	c.V[x] = byte(z & 0x0f)
	log.Printf("%x ADD V%d, V%d", c.PC, x, y)
}
