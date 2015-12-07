package main

// Memory interface describes the fetch and write methods for the VM
type Memory interface {
	Fetch(uint16) byte
	Write(uint16, byte)
}

type StdMemory []byte

// NewBlock returns a new Standard Memory block of at least sz bytes
// If sz is zero then NewBlock returns the maximum value of a uint16
func NewBlock(sz uint16) StdMemory {
	if sz == 0 {
		sz = 0xffff
	}
	return make(StdMemory, sz)
}

func (s StdMemory) Fetch(addr uint16) byte {
	if uint16(len(s)) <= addr {
		panic("segfault: address out of bounds")
	}
	return s[addr]
}

func (s StdMemory) Write(addr uint16, data byte) {
	if uint16(len(s)) <= addr {
		panic("segfault: address out of bounds")
	}
	s[addr] = data
}

func (s StdMemory) WriteBlock(addr uint16, data []byte) {
	if uint16(len(s)) <= addr {
		panic("segfault: address out of bounds")
	}
	copy(s[addr:], data)
}
