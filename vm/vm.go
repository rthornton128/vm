package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	vm "github.com/rthornton128/vm/lib"
)

type CPU struct {
	ar   uint16 // address register
	dr   byte   // data register
	ir   byte   // instruction register
	pc   uint16 // program counter
	sp   uint16 // stack pointer
	tr   byte
	ac   byte // accumulator
	b    byte // register b
	c    byte // register c
	zero bool
	mem  []byte
	end  uint16
}

func (c *CPU) init(ep uint16, prog []byte) {
	c.ac = 0 // redundant but extra assurrance
	c.pc = ep
	c.sp = uint16(len(prog))
	c.mem = make([]byte, 0xffff) // 65536 bytes
	c.zero = true

	// copy program in
	copy(c.mem, prog)

	// set return address on stack to invalid address
	c.mem[c.sp] = 0xff   // lsb
	c.mem[c.sp+1] = 0xff // msb
	c.sp += 2
}

func (c *CPU) decode() {
	switch vm.Opcode(c.ir & 0x3f) {
	case vm.NOP:
	case vm.CALL:
		// load address to jump to
		c.ar = c.pc
		c.pc++
		c.dr = c.mem[c.ar]
		c.ar = c.pc
		c.pc++
		c.tr = c.dr
		c.dr = c.mem[c.ar]
	case vm.JMP, vm.JPZ, vm.JNZ:
		c.dr = c.mem[c.pc]
		c.ar = uint16(c.dr) << 8
		c.pc++
		c.dr = c.mem[c.pc]
		c.ar |= uint16(c.dr)
		c.pc++
	case vm.RET:
		c.sp--
		c.ar = c.sp
		c.dr = c.mem[c.ar]
		c.sp--
		c.ar = c.sp
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
		c.ar = c.sp
		c.sp++
		c.mem[c.ar] = uint8(c.pc) // >> 8)
		c.ar = c.sp
		c.sp++
		c.mem[c.ar] = uint8(c.pc >> 8)
		c.pc = uint16(c.tr) << 8
		c.pc |= uint16(c.dr)
	case vm.RET:
		c.pc = uint16(c.dr) << 8
		c.dr = c.mem[c.ar]
		c.pc |= uint16(c.dr)
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
	c.pc++
	c.dr = c.mem[c.ar] // fetch instruction into data register
	c.ir = c.dr
	//fmt.Println("c.ar:", c.ar, "c.ir:", c.ir)
}

func (c *CPU) run() {
	for {
		c.fetch()
		c.decode()
		c.exec()
		if c.pc == 0xffff {
			break
		}
	}
}

func main() {
	flag.Parse()

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	p, err := vm.ScanProgram(b)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(p)
	//fmt.Println(p.SecTab[vm.TEXT])
	cpu := new(CPU)
	cpu.init(p.Entry, p.SecTab[vm.TEXT])
	cpu.run()
	log.Println("exit with result:", cpu.ac)
}
