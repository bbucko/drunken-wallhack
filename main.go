package main

import (
	"io/ioutil"
	"log"

	"time"

	"github.com/bbucko/drunken-wallhack/chip8"
)

func main() {
	data, err := ioutil.ReadFile("roms/MAZE")
	if err != nil {
		log.Fatal(err)
	}

	c := chip8.New(data)

	for {
		err := c.Step()
		time.Sleep(100 * time.Millisecond)
		if err != nil {
			log.Fatal("", err)
		}
	}

}
