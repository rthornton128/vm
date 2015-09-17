package vm_test

import (
	"bytes"
	"testing"

	vm "github.com/rthornton128/vm/lib"
)

/*
func TestObject(t *testing.T) {
	o := &vm.Object{
		Entry: 0x03,
		SecTab: vm.SectionTable{
			vm.OSection{
				Name: "text",
				Data: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa},
			},
		},
		RelocTab: vm.RelocateTable{0x2: 0x6},
		SymTab: vm.SymbolTable{
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
	if len(ob.SymTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, s := range ob.SymTab {
		if s != (o.SymTab)[i] {
			t.Log("expected:", (ob.SymTab)[i], "got:", s)
			t.FailNow()
		}
	}
	if len(ob.RelocTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, r := range ob.RelocTab {
		if r != o.RelocTab[i] {
			t.Log("expected:", o.RelocTab[i], "got:", r)
			t.FailNow()
		}
	}
	if len(ob.SecTab) < 1 {
		t.Fatal("failed to scan symbol table")
	}
	for i, s := range ob.SecTab {
		if s.Name != (o.SecTab)[i].Name ||
			!bytes.Equal(s.Data, (o.SecTab)[i].Data) {
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

func TestRelocateTable(t *testing.T) {
	rt := vm.RelocateTable{
		0x42: 0xabcd,
		0xff: 0x1234,
	}
	b := rt.Bytes()
	//expect := []byte{0x42, 0xab, 0xcd, 0xff, 0x12, 0x34}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	rel := vm.ScanRelocateTable(b, uint16(len(b)))
	for i, r := range rel {
		if r != rt[i] {
			t.Log("expected:", rt, "got:", rel)
			t.FailNow()
		}
	}
}

func TestRelocateTableMerge(t *testing.T) {
	rt1 := vm.RelocateTable{0x0: 0x0000, 0x1: 0x0001}
	rt2 := vm.RelocateTable{0x2: 0x0002, 0x3: 0x0003}
	expect := []byte{0x6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1,
		0x2, 0x0, 0x2, 0x3, 0x0, 0x3}

	err := rt1.Merge(rt2)

	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(rt1.Bytes(), expect) {
		t.Fatal("expected:", expect, "got:", rt1.Bytes())
	}
}
*/
func TestSectionTable(t *testing.T) {
	st := vm.SectionTable{
		vm.TEXT: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
		vm.DATA: []byte{0xa, 0xb, 0xc, 0xd, 0xe, 0xf},
	}
	b := st.Bytes()
	expect := []byte{0x2,
		byte(vm.TEXT), 0x0, 0xb, 0x0, 0x5,
		byte(vm.DATA), 0x0, 0x10, 0x0, 0x6,
		0x1, 0x2, 0x3, 0x4, 0x5,
		0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
	}
	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	o := vm.NewObject()
	o.ScanSectionTable(b)
	for i, sec := range o.SecTab {
		if !bytes.Equal(st[i], sec) {
			t.Log("expected:", st[i], "got:", sec)
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
	st := vm.SymbolTable{
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

	o := vm.NewObject()
	o.ScanSymbolTable(b)
	for i, s := range o.SymTab {
		if s != st[i] {
			t.Log("expected:", st, "got:", o.SymTab)
			t.FailNow()
		}
	}
}
