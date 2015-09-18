package vm

import (
	"bytes"
	"errors"
	"fmt"
)

var MagicNumber = []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}

type Object struct {
	Entry    uint16
	RelAddr  uint16
	RelSize  uint16
	SecAddr  uint16
	SecSize  uint16
	SymAddr  uint16
	SymSize  uint16
	SecTab   SectionTable
	RelocTab RelocateTable
	SymTab   SymbolTable
}

func NewObject() *Object {
	return &Object{
		SecTab:   make(SectionTable, section_max),
		RelocTab: make(RelocateTable, 0),
		SymTab:   make(SymbolTable, 0),
	}
}

func ScanObject(b []byte) (*Object, error) {
	o := NewObject()

	// TODO Not enough...see below
	if len(b) < 12 {
		return o, errors.New("not enough bytes")
	}

	if !bytes.Equal(b[:8], MagicNumber) {
		return o, errors.New("bad magic number, not vm object")
	}

	o.Entry = toAddress(b[8:10])

	o.RelAddr = toAddress(b[10:12])
	o.RelSize = toAddress(b[12:14])

	o.SymAddr = toAddress(b[14:16])
	o.SymSize = toAddress(b[16:18])

	o.SecAddr = toAddress(b[18:20])
	o.SecSize = toAddress(b[20:22])

	o.ScanRelocateTable(b[o.RelAddr : o.RelAddr+o.RelSize])
	o.ScanSymbolTable(b[o.SymAddr : o.SymAddr+o.SymSize])
	o.ScanSectionTable(b[o.SecAddr : o.SecAddr+o.SecSize])

	return o, nil
}

func (o *Object) Bytes() []byte {
	b := make([]byte, 0)

	// magic #
	b = append(b, MagicNumber...)

	// entry point
	b = append(b, toBytes(o.Entry)...)

	i := uint16(len(MagicNumber)) + 14

	// relocation table size and addr
	b = append(b, toBytes(i)...)
	b = append(b, toBytes(o.RelocTab.Size())...)
	i += o.RelocTab.Size()

	// symbol table size and addr
	b = append(b, toBytes(i)...)
	b = append(b, toBytes(o.SymTab.Size())...)
	i += o.SymTab.Size()

	// section table size and addr
	b = append(b, toBytes(i)...)
	b = append(b, toBytes(o.SecTab.Size())...)
	i += o.SecTab.Size()

	// tables
	b = append(b, o.RelocTab.Bytes()...)
	b = append(b, o.SymTab.Bytes()...)
	b = append(b, o.SecTab.Bytes()...)

	return b
}

func (o *Object) Merge(objs ...*Object) error {
	for _, ob := range objs {
		if err := o.MergeRelocates(ob.RelocTab); err != nil {
			return err
		}
		if err := o.MergeSymbols(ob.SymTab); err != nil {
			return err
		}
		// merge sections last so relocations can happen
		if err := o.MergeSections(ob.SecTab); err != nil {
			return err
		}
	}

	// TODO temporary hack
	for _, s := range o.SymTab {
		if s.Name == "main" {
			o.Entry = s.Addr
		}
	}
	return nil
}

type SecType byte

const (
	TEXT SecType = iota
	DATA
	section_max
)

var sections = []string{
	DATA: "data",
	TEXT: "text",
}

func LookupSectionName(name string) (byte, error) {
	for i, s := range sections {
		if name == s {
			return byte(i), nil
		}
	}
	return 0, fmt.Errorf("invalid section name: %s", name)
}

// Sectiontable is a segment of data that may represent different parts of the
// program. The text section contains binary opcodes for the virtual machine
// to execute.
type SectionTable [][]byte

func (o *Object) ScanSectionTable(b []byte) error {
	// section table starts with the number of sections that should be scanned
	nsec := int(b[0])

	// it then contains a table with format: type, address, length
	for i, j := 0, 1; i < nsec; i++ {
		//fmt.Println(i, j)
		t := b[j]
		if cap(o.SecTab[t]) > 0 {
			return fmt.Errorf("duplicate section:", sections[t])
		}
		addr := toAddress(b[j+1 : j+3])
		ln := toAddress(b[j+3 : j+5])

		o.SecTab[t] = make([]byte, ln)
		copy(o.SecTab[t], b[addr:addr+ln])
		j += 5
	}
	return nil
}

func (o *Object) setSection(sec SecType, data []byte) error {
	if sec >= section_max {
		return fmt.Errorf("invalid section: %d", sec)
	}

	if cap(o.SecTab[sec]) > 0 {
		return fmt.Errorf("duplicate section not allowed: %s", sec)
	}

	o.SecTab[sec] = make([]byte, len(data))
	copy(o.SecTab[sec], data)
	//fmt.Println(o.SecTab)

	return nil
}

func (st SectionTable) Bytes() []byte {
	sz := 1
	// TODO likely a better way
	var n byte
	for i := range st {
		if cap(st[i]) > 0 {
			n++
			sz += 5 + len(st[i])
		}
	}
	b := make([]byte, sz)
	b[0] = byte(n)
	i, j := 1, int(1+(n*5))
	for t, sec := range st {
		if cap(st[t]) > 0 {
			b[i] = byte(t)
			copy(b[i+1:], toBytes(uint16(j)))
			copy(b[i+3:], toBytes(uint16(len(sec))))
			copy(b[j:], sec)
			i, j = i+5, j+len(sec)
		}
	}
	return b
}

func (o *Object) MergeSections(other SectionTable) error {
	// TODO should it return error? will one occur?
	for sec, data := range other {
		if cap(o.SecTab[sec]) > 0 {
			// TODO intentionally wrong; needs relocations and symbol table update
			o.SecTab[sec] = append(o.SecTab[sec], data...)
		} else {
			o.SecTab[sec] = make([]byte, len(data))
			copy(o.SecTab[sec], data)
		}
	}
	return nil
}

func (st SectionTable) Size() uint16 {
	sz := 0
	for _, v := range st {
		if len(v) > 0 {
			sz += len(v) + 5
		}
	}
	return uint16(sz)
}

// Relocate holds the offset of an address within the text section of an
// object. It also contains an index into the symbol table. A relocate
// object is used by the linker to adjust the location of symbols in
// memory
// RelocateTable is a list of relocatable objects
type RelocateTable []RelocAddr

type RelocAddr struct {
	index  byte
	offset uint16
}

func (o *Object) ScanRelocateTable(b []byte) {
	for i := 0; i < len(b); i += 3 {
		o.AddRelocate(b[0], toAddress(b[1:3]))
		//o.RelocTab = append(o.RelocTab,
		//RelocAddr{index: b[0], offset: toAddress(b[1:3])})
	}
	return
}

func (o *Object) AddRelocate(index byte, offset uint16) error {
	for _, r := range o.RelocTab {
		if r.index == index {
			return fmt.Errorf("duplicate symbolic reference: %d", index)
		}
	}
	o.RelocTab = append(o.RelocTab, RelocAddr{index, offset})
	return nil
}

func (rt RelocateTable) Bytes() []byte {
	b := make([]byte, 0)
	for _, r := range rt {
		b = append(b, r.index)
		b = append(b, toBytes(r.offset)...)
	}
	return b
}

func (o *Object) MergeRelocates(other RelocateTable) error {
	for _, r := range other {
		if err := o.AddRelocate(r.index, r.offset); err != nil {
			return err
		}
	}
	return nil
}

func (rt RelocateTable) Size() uint16 {
	return uint16(len(rt) * 3)
}

// Symbol represent an addressable location associated with a label.
// Function and variable names are examples
type Symbol struct {
	Name string
	Addr uint16
}

func ScanSymbol(b []byte) Symbol {
	sz := b[2]
	return Symbol{
		Addr: toAddress(b[:2]),
		Name: string(b[3 : 3+sz]),
	}
}

func (s Symbol) Bytes() []byte {
	b := toBytes(s.Addr)
	b = append(b, uint8(len(s.Name)))
	b = append(b, []byte(s.Name)...)

	return b
}

func (s Symbol) Size() uint16 {
	return uint16(len(s.Name) + 3)
}

// SymbolTable is a list of all Symbols found in the object/program
type SymbolTable []Symbol

func (o *Object) ScanSymbolTable(b []byte) {
	for i := 0; i < len(b); {
		s := ScanSymbol(b[i:])
		o.SymTab = append(o.SymTab, s)
		i += int(s.Size())
	}
	return
}

func (o *Object) AddSymbol(name string, addr uint16) error {
	for _, sym := range o.SymTab {
		if sym.Name == name {
			return fmt.Errorf("duplicate name: %s", name)
		}
	}
	o.SymTab = append(o.SymTab, Symbol{Addr: addr, Name: name})
	return nil
}

func (st SymbolTable) Bytes() []byte {
	b := make([]byte, 0)
	for _, s := range st {
		b = append(b, s.Bytes()...)
	}
	return b
}

func (o *Object) MergeSymbols(other SymbolTable) error {
	for _, sym := range other {
		if err := o.AddSymbol(sym.Name, sym.Addr); err != nil {
			return err
		}
	}
	return nil
}

func (o *Object) LookupSymbolIndex(name string) byte {
	for i, s := range o.SymTab {
		if name == s.Name {
			return byte(i)
		}
	}
	return 255
}

func (st SymbolTable) Lookup(name string) (Symbol, bool) {
	for _, s := range st {
		if name == s.Name {
			return s, true
		}
	}
	return Symbol{}, false
}

func (st SymbolTable) Size() uint16 {
	var sz uint16
	for _, s := range st {
		sz += s.Size()
	}
	return sz
}
