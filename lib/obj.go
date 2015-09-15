package vm

import (
	"bytes"
	"errors"
	"fmt"
)

var MagicNumber = []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}

type Object struct {
	Entry    uint16
	SecTab   *SectionTable
	RelocTab *RelocateTable
	SymTab   *SymbolTable
}

func NewObject() *Object {
	sec := make(SectionTable, 0)
	rel := make(RelocateTable, 0)
	sym := make(SymbolTable, 0)

	return &Object{
		SecTab:   &sec,
		RelocTab: &rel,
		SymTab:   &sym,
	}
}

func ScanObject(b []byte) (*Object, error) {
	//var o Object
	o := NewObject()

	// TODO Not enough...see below
	if len(b) < 12 {
		return o, errors.New("not enough bytes")
	}

	if !bytes.Equal(b[:8], MagicNumber) {
		return o, errors.New("bad magic number, not vm object")
	}

	o.Entry = toAddress(b[8:10])

	// The basic header/object should actually use addresses and sizes of
	// each section so that certain sections can be skipped if they do not
	// exist
	sz := toAddress(b[10:12])
	o.RelocTab = ScanRelocateTable(b[12:], sz)

	last := 12 + sz
	sz = toAddress(b[last : last+2])
	o.SymTab = ScanSymbolTable(b[last+2:], sz)

	last += sz + 2
	sz = toAddress(b[last : last+2])
	o.SecTab = ScanSectionTable(b[last+2:], sz)

	return o, nil
}

func (o *Object) Bytes() []byte {
	b := make([]byte, 0)

	// magic #
	b = append(b, MagicNumber...)

	// entry point
	b = append(b, toBytes(o.Entry)...)

	// relocation table
	b = append(b, toBytes(o.RelocTab.Size())...)
	b = append(b, o.RelocTab.Bytes()...)

	// symbol table
	b = append(b, toBytes(o.SymTab.Size())...)
	b = append(b, o.SymTab.Bytes()...)

	// section table
	b = append(b, toBytes(o.SecTab.Size())...)
	b = append(b, o.SecTab.Bytes()...)

	return b
}

func (o *Object) Merge(objs ...*Object) error {
	for _, ob := range objs {
		if err := o.RelocTab.Merge(ob.RelocTab); err != nil {
			return err
		}
		if err := o.SymTab.Merge(ob.SymTab); err != nil {
			return err
		}
		// merge sections last so relocations can happen
		if err := o.SecTab.Merge(ob.SecTab); err != nil {
			return err
		}
	}

	// TODO temporary hack
	for _, s := range *o.SymTab {
		if s.Name == "main" {
			o.Entry = s.Addr
		}
	}
	return nil
}

// Section is a segment of data that may represent different parts of the
// program. The text section contains binary opcodes for the virtual machine
// to execute.
type OSection struct {
	Name string
	Data []byte
}

func ScanSection(b []byte) OSection {
	strlen := b[0]
	datastart := uint16(strlen) + 3
	datalen := toAddress(b[strlen+1 : strlen+3])
	return OSection{
		Name: string(b[1 : strlen+1]),
		Data: b[datastart : datastart+datalen],
	}
}

func (os OSection) Bytes() []byte {
	b := make([]byte, 0)

	b = append(b, byte(len(os.Name))) // max len 255
	b = append(b, []byte(os.Name)...)

	b = append(b, toBytes(uint16(len(os.Data)))...)
	b = append(b, os.Data...)

	return b
}

func (os OSection) Merge(other OSection) error {
	if os.Name != other.Name {
		return fmt.Errorf("merge: sections do not match: %s vs %s",
			os.Name, other.Name)
	}
	// TODO temporary stub, needs to do relocations and lookups
	os.Data = append(os.Data, other.Data...)
	return nil
}

func (os OSection) Size() uint16 {
	return uint16(len(os.Name) + len(os.Data) + 3)
}

type SectionTable []OSection

func ScanSectionTable(b []byte, n uint16) *SectionTable {
	st := make(SectionTable, 0)
	for i := uint16(0); i < n; {
		s := ScanSection(b[i:])
		st = append(st, s)
		i += s.Size()
	}
	return &st
}

func (st *SectionTable) Add(os OSection) error {
	for _, sec := range *st {
		if sec.Name == os.Name {
			return fmt.Errorf("duplicate section not allowed: %s", sec.Name)
		}
	}
	*st = append(*st, os)
	fmt.Println(*st)
	return nil
}

func (st *SectionTable) Bytes() []byte {
	b := make([]byte, 0)
	for _, s := range *st {
		b = append(b, s.Bytes()...)
	}
	return b
}

func (st *SectionTable) Merge(other *SectionTable) error {
	for _, sec := range *other {
		fmt.Println("sect add:", sec)
		if err := st.MergeSection(sec); err != nil {
			return err
		}
	}
	return nil
}

func (st *SectionTable) MergeSection(other OSection) error {
	for _, sec := range *st {
		fmt.Println("sect merge:", sec)
		if sec.Name == other.Name {
			if err := sec.Merge(other); err != nil {
				return err
			}
		}
	}
	*st = append(*st, other)
	return nil
}

func (st *SectionTable) Size() uint16 {
	var sz uint16
	for _, s := range *st {
		sz += s.Size()
	}
	return sz
}

// Relocate holds the offset of an address within the text section of an
// object. It also contains an index into the symbol table. A relocate
// object is used by the linker to adjust the location of symbols in
// memory
type Relocate struct {
	Offset   uint16
	SymIndex byte
}

func ScanRelocate(b []byte) Relocate {
	return Relocate{SymIndex: b[0], Offset: toAddress(b[1:3])}
}

func (r Relocate) Bytes() []byte {
	b := toBytes(r.Offset)
	return []byte{r.SymIndex, b[0], b[1]}
}

// RelocateTable is a list of relocatable objects
type RelocateTable []Relocate

func ScanRelocateTable(b []byte, n uint16) *RelocateTable {
	rt := make(RelocateTable, 0)
	for i := uint16(0); i < n; i += 3 {
		rt = append(rt, ScanRelocate(b[i:]))
	}
	return &rt
}

func (rt *RelocateTable) Add(r Relocate) error {
	for _, reloc := range *rt {
		if reloc.SymIndex == r.SymIndex {
			return fmt.Errorf("duplicate symbolic reference: %d", r.SymIndex)
		}
	}
	*rt = append(*rt, r)
	fmt.Println(*rt)
	return nil
}

func (rt *RelocateTable) Bytes() []byte {
	b := make([]byte, 0)
	for _, r := range *rt {
		b = append(b, r.Bytes()...)
	}
	return b
}

func (rt *RelocateTable) Merge(other *RelocateTable) error {
	for _, reloc := range *other {
		fmt.Println("rt add:", *other)
		if err := rt.Add(reloc); err != nil {
			return err
		}
	}
	fmt.Println(*rt)
	return nil
}

func (rt *RelocateTable) Size() uint16 {
	return uint16(len(*rt) * 3)
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

func ScanSymbolTable(b []byte, n uint16) *SymbolTable {
	st := make(SymbolTable, 0)
	for i := uint16(0); i < n; {
		s := ScanSymbol(b[i:])
		st = append(st, s)
		i += s.Size()
	}
	return &st
}

func (st *SymbolTable) Add(s Symbol) error {
	for _, sym := range *st {
		if sym.Name == s.Name {
			return fmt.Errorf("duplicate name: %s", sym.Name)
		}
	}
	*st = append(*st, s)
	return nil
}

func (st *SymbolTable) Bytes() []byte {
	b := make([]byte, 0)
	for _, s := range *st {
		b = append(b, s.Bytes()...)
	}
	return b
}

func (st *SymbolTable) Merge(other *SymbolTable) error {
	for _, sym := range *other {
		fmt.Println("st add:", *other)
		if err := st.Add(sym); err != nil {
			return err
		}
	}
	return nil
}

func (st *SymbolTable) Lookup(name string) (Symbol, bool) {
	for _, s := range *st {
		if name == s.Name {
			return s, true
		}
	}
	return Symbol{}, false
}

func (st *SymbolTable) Size() uint16 {
	var sz uint16
	for _, s := range *st {
		sz += s.Size()
	}
	return sz
}
