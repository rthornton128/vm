package vm

import (
	"errors"
	"fmt"
)

type Address uint16

func (a Address) String() string {
	return fmt.Sprintf("%x", a)
}

type Opcode byte

const (
	NOP Opcode = iota // no operation

	/* Branching */
	JMP
	JPZ
	JNZ
	CALL
	RET

	/* Register */
	MOV
	MVR
	MVI
	CLA
	CLR

	/* Stack */
	POP
	PUSH

	/* Arithmetic */
	ADD // add
	DIV // divide
	INC
	MUL // multiply
	SHL
	SHR
	SUB // subtract

	/* Logical */
	AND
	OR
)

var opcodes = map[Opcode]string{
	NOP:  "nop",
	JMP:  "jmp",
	JPZ:  "jpz", // jump if zero
	JNZ:  "jnz", // jump if not zero
	CALL: "call",
	RET:  "ret",
	MOV:  "mov",
	MVR:  "mvr",
	MVI:  "mvi",
	CLA:  "cla",
	CLR:  "clr",
	POP:  "pop",
	PUSH: "push",
	ADD:  "add",
	DIV:  "div",
	INC:  "inc",
	MUL:  "mul",
	SHL:  "shl",
	SHR:  "shr",
	SUB:  "sub",
	AND:  "and",
	OR:   "or",
}

func (o Opcode) String() string {
	return opcodes[o]
}

func LookupOpcode(s string) (Opcode, error) {
	for k, v := range opcodes {
		if s == v {
			return k, nil
		}
	}
	return 0, errors.New("invalid instruction: " + s)
}

type Register byte

const (
	REGB Register = 0x00
	REGC Register = 0x80
)

var registers = map[Register]string{
	REGB: "b",
	REGC: "c",
}

func (r Register) String() string {
	return registers[r]
}

func LookupRegister(s string) (Register, error) {
	for k, v := range registers {
		if s == v {
			return k, nil
		}
	}
	return 0, errors.New("invalid register name: " + s)
}
