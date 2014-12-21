package main

import (
	"io/ioutil"
	"log"

	"github.com/bbucko/drunken-wallhack/chip8"
)

func main() {
	data, err := ioutil.ReadFile("roms/zero")
	if err != nil {
		log.Fatal(err)
	}

	c := chip8.New()
	c.Load(data)

	for i := 0; i < 10; i++ {
		c.Step()
	}

}
