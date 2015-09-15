package vm_test

import (
	"bytes"
	"testing"

	vm "github.com/rthornton128/vm/lib"
)

func TestObject(t *testing.T) {
	o := &vm.Object{
		Entry: 0x03,
		SecTab: &vm.SectionTable{
			vm.OSection{
				Name: "text",
				Data: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa},
			},
		},
		RelocTab: &vm.RelocateTable{
			vm.Relocate{
				Offset:   0x6,
				SymIndex: 0x2,
			},
		},
		SymTab: &vm.SymbolTable{
			vm.Symbol{
				Name: "fn",
				Addr: 0x0,
			},
			vm.Symbol{
				Name: "main",
				Addr: 0x3,
			},
		},
	}
	b := o.Bytes()
	expect := []byte{
		0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf, // magic #
		0x0, 0x3, // entry pt
		0x0, 0x3, // reloc len
		0x2, 0x0, 0x6, // reloc
		0x0, 0xc, // symtab len
		0x0, 0x0, 0x2, 'f', 'n', // symbol 1
		0x0, 0x3, 0x4, 'm', 'a', 'i', 'n', // symbol 2
		0x0, 0x11, // section len
		0x4, 't', 'e', 'x', 't',
		0x0, 0xa, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa,
	}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	ob, err := vm.ScanObject(b)
	if err != nil {
		t.Fatal(err)
	}
	if ob.Entry != o.Entry {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}
	if len(*ob.SymTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, s := range *ob.SymTab {
		if s != (*o.SymTab)[i] {
			t.Log("expected:", (*ob.SymTab)[i], "got:", s)
			t.FailNow()
		}
	}
	if len(*ob.RelocTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, r := range *ob.RelocTab {
		if r != (*o.RelocTab)[i] {
			t.Log("expected:", (*o.RelocTab)[i], "got:", r)
			t.FailNow()
		}
	}
	if len(*ob.SecTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, s := range *ob.SecTab {
		if s.Name != (*o.SecTab)[i].Name ||
			!bytes.Equal(s.Data, (*o.SecTab)[i].Data) {
			t.Log("expected:", ob.SecTab, "got:", o.SecTab)
			t.FailNow()
		}
	}

	_, err = vm.ScanObject(b[:10])
	if err == nil {
		t.Fatal("expected error, got none")
	}

	b[0] = 42
	_, err = vm.ScanObject(b)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestRelocate(t *testing.T) {
	r := vm.Relocate{Offset: 0x24, SymIndex: 0x1}
	b := r.Bytes()
	expect := []byte{0x1, 0x0, 0x24}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	rlc := vm.ScanRelocate(b)
	if rlc != r {
		t.Log("expected:", r, "got:", rlc)
		t.FailNow()
	}
}

func TestRelocateTable(t *testing.T) {
	rt := &vm.RelocateTable{
		vm.Relocate{Offset: uint16(0xabcd), SymIndex: 0x42},
		vm.Relocate{Offset: uint16(0x1234), SymIndex: 0xff},
	}
	b := rt.Bytes()
	expect := []byte{0x42, 0xab, 0xcd, 0xff, 0x12, 0x34}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	rel := vm.ScanRelocateTable(b, uint16(len(b)))
	for i, r := range *rel {
		if r != (*rt)[i] {
			t.Log("expected:", rt, "got:", rel)
			t.FailNow()
		}
	}
}

func TestSection(t *testing.T) {
	s := vm.OSection{Name: "text", Data: []byte{0x1, 0x2, 0x3, 0x4, 0x5}}
	b := s.Bytes()
	expect := []byte{0x4, 't', 'e', 'x', 't', 0x0, 0x5, 0x1, 0x2, 0x3, 0x4, 0x5}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	sec := vm.ScanSection(b)
	if sec.Name != s.Name || !bytes.Equal(sec.Data, s.Data) {
		t.Log("expected:", s, "got:", sec)
		t.FailNow()
	}
}

func TestSectionTable(t *testing.T) {
	st := &vm.SectionTable{
		vm.OSection{Name: "text", Data: []byte{0x1, 0x2, 0x3, 0x4, 0x5}},
		vm.OSection{Name: "data", Data: []byte{0xa, 0xb, 0xc, 0xd, 0xe, 0xf}},
	}
	b := st.Bytes()
	expect := []byte{
		0x4, 't', 'e', 'x', 't', 0x0, 0x5, 0x1, 0x2, 0x3, 0x4, 0x5,
		0x4, 'd', 'a', 't', 'a', 0x0, 0x6, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
	}
	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	sec := vm.ScanSectionTable(b, uint16(len(b)))
	for i, s := range *sec {
		if (*st)[i].Name != s.Name || !bytes.Equal((*st)[i].Data, s.Data) {
			t.Log("expected:", st, "got:", sec)
			t.FailNow()
		}
	}
}

func TestSymbol(t *testing.T) {
	s := vm.Symbol{Addr: uint16(0x3), Name: "main"}
	b := s.Bytes()
	expect := []byte{0x0, 0x3, 0x4, 'm', 'a', 'i', 'n'}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	sym := vm.ScanSymbol(b)
	if sym != s {
		t.Log("expected:", s, "got:", sym)
		t.FailNow()
	}
}

func TestSymbolTable(t *testing.T) {
	st := &vm.SymbolTable{
		vm.Symbol{Addr: uint16(0xabcd), Name: "foo"},
		vm.Symbol{Addr: uint16(0x1234), Name: "bar"},
	}
	b := st.Bytes()
	expect := []byte{
		0xab, 0xcd, 0x3, 'f', 'o', 'o', 0x12, 0x34, 0x3, 'b', 'a', 'r',
	}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	sym := vm.ScanSymbolTable(b, uint16(len(b)))
	for i, s := range *sym {
		if s != (*st)[i] {
			t.Log("expected:", st, "got:", sym)
			t.FailNow()
		}
	}
}
