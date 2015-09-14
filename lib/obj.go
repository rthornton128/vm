package vm

var MagicNumber = []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf}

type Object struct {
	Entry uint16
	//Size uint16
	Sections []ObjectSection
	RelocTab RelocateTable
	SymTab   SymbolTable
}

func (o Object) Bytes() []byte {
	b := make([]byte, 0)

	// magic #
	b = append(b, MagicNumber...)

	// entry
	u, l := toBytes(o.Entry)
	b = append(b, u, l)

	// sections
	for _, s := range o.Sections {
		b = append(b, s.Bytes()...)
	}

	// Relocatable Symbols
	for _, r := range o.RelocTab {
		b = append(b, r.Bytes()...)
	}

	// sections
	u, l = toBytes(o.SymTab.sz)
	b = append(b, u, l)
	for _, s := range o.SymTab.table {
		b = append(b, s.Bytes()...)
	}

	return b
}

func (o Object) Merge(objs ...Object) {
	for _ = range objs {
	}
}

type ObjectSection struct {
	Name string
	Data []byte
	Size uint16
}

func (os ObjectSection) Bytes() []byte {
	b := make([]byte, 0)
	u, l := toBytes(uint16(len(os.Name)))
	b = append(b, u, l)
	b = append(b, []byte(os.Name)...)
	u, l = toBytes(uint16(len(os.Data)))
	b = append(b, u, l)
	b = append(b, os.Data...)

	return b
}

type Relocate struct {
	Offset   uint16
	SymIndex byte
}

func (r Relocate) Bytes() []byte {
	u, l := toBytes(r.Offset)
	return []byte{r.SymIndex, u, l}
}

type RelocateTable []Relocate

type Symbol struct {
	Name string
	Type byte
	Addr uint16
}

func (s Symbol) Bytes() []byte {
	u, l := toBytes(s.Addr)
	b := []byte{u, l, s.Type, uint8(len(s.Name))}
	b = append(b, []byte(s.Name)...)

	return b
}

func (s Symbol) Size() uint16 {
	return uint16(len(s.Name)) + 3
}

type SymbolTable struct {
	table []Symbol
	sz    uint16
}
