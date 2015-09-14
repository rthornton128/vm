package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	vm "github.com/rthornton128/vm/lib"
)

type instruction struct {
	i    vm.Instruction
	sz   int
	rega vm.Register
	regb vm.Register
	imm  byte
	addr uint16
}

type CPU struct {
	ar   uint16 // address register
	dr   byte   // data register
	ir   byte   // instruction register
	pc   uint16 // program counter
	sp   uint16 // stack pointer
	ac   byte   // accumulator
	b    byte   // register b
	c    byte   // register c
	zero bool
	mem  []byte
	end  uint16
}

func newCPU(mem []byte) *CPU {
	return &CPU{}
}

func (c *CPU) decode() {
	switch vm.Opcode(c.ir & 0x3f) {
	case vm.NOP:
	case vm.CALL:
		c.ar = c.sp
		c.mem[c.ar] = uint8(c.pc >> 8)
		c.sp++
		c.ar = c.sp
		c.mem[c.ar] = uint8(c.pc)
		c.sp++
		fallthrough
	case vm.JMP, vm.JPZ, vm.JNZ:
		c.dr = c.mem[c.pc]
		c.ar = uint16(c.dr) << 8
		c.pc++
		c.dr = c.mem[c.pc]
		c.ar |= uint16(c.dr)
		c.pc++
	case vm.RET:
	case vm.POP:
		c.sp--
		c.ar = c.sp
	case vm.PUSH:
		c.dr = c.ac
		c.ar = c.sp
		c.sp++
	case vm.MVR:
		c.dr = c.ac
	case vm.MVI:
		c.ar = c.pc
		c.dr = c.mem[c.ar]
		c.pc++
	case vm.MOV, vm.ADD, vm.DIV, vm.MUL, vm.SHL, vm.SHR, vm.SUB, vm.AND, vm.OR:
		switch vm.Register(c.ir & 0xc0) {
		case vm.REGB:
			c.dr = c.b
		case vm.REGC:
			c.dr = c.c
		}
	case vm.INC:
		c.dr = 1
	}
}

func (c *CPU) exec() {
	switch vm.Opcode(c.ir & 0x3f) {
	case vm.NOP:
	case vm.JMP:
		c.pc = c.ar
	case vm.JPZ:
		if c.zero {
			c.pc = c.ar
		}
	case vm.JNZ:
		if !c.zero {
			c.pc = c.ar
		}
	case vm.CALL:
	case vm.RET:
	case vm.MOV:
		c.ac = c.dr
	case vm.MVR:
		switch vm.Register(c.ir & 0xc0) {
		case vm.REGB:
			c.b = c.dr
		case vm.REGC:
			c.c = c.dr
		}
	case vm.MVI:
		c.ac = c.dr
	case vm.CLA:
		c.ac = 0
		c.zero = true
	case vm.CLR:
		c.ac = 0
	case vm.POP:
		c.dr = c.mem[c.ar]
		c.ac = c.dr
	case vm.PUSH:
		c.mem[c.ar] = c.dr
	case vm.ADD, vm.INC:
		c.ac += c.dr
	case vm.DIV:
		c.ac /= c.dr
	case vm.MUL:
		c.ac *= c.dr
	case vm.SHL:
		c.ac <<= c.dr
	case vm.SHR:
		c.ac >>= c.dr
	case vm.SUB:
		c.ac -= c.dr
	case vm.AND:
		c.ac &= c.dr
	case vm.OR:
		c.ac |= c.dr
	}
}

func (c *CPU) fetch() {
	c.ar = c.pc
	c.dr = c.mem[c.ar] // fetch instruction into data register
	c.pc++
	c.ir = c.dr
}

func (c *CPU) run() {
	c.fetch()
}

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	magic := make([]byte, 8)

	// verify file is a valid virtual machine object
	n, err := f.Read(magic)
	if err != nil {
		log.Fatal(err)
	}
	if n != 8 {
		log.Fatal("failed to read 8 bytes, not vm object file")
	}
	if !bytes.Equal([]byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}, magic) {
		log.Fatal("not a vm object file")
	}

	// load entry point
	ep := make([]byte, 2)
	n, err = f.Read(ep)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("failed to read entry point")
	}

	// load program into memory
	mem := make([]byte, 0xffff) // 65536 byte array
	_, err = f.Read(mem)
	if err != nil {
		log.Fatal(err)
	}

	cpu := newCPU(mem)
	cpu.run()
}
