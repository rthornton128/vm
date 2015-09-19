package vm_test

import (
	"bytes"
	"log"
	"testing"

	vm "github.com/rthornton128/vm/lib"
)

func TestRelocateTable(t *testing.T) {
	o := vm.NewObject()
	o.AddRelocate(0x42, 0xabcd)
	o.AddRelocate(0xff, 0x1234)
	b := o.RelocTab.Bytes()
	expect := []byte{0x42, 0xab, 0xcd, 0xff, 0x12, 0x34}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	o2 := vm.NewObject()
	o2.ScanRelocateTable(b)
	for i, r := range o2.RelocTab {
		if r != o.RelocTab[i] {
			t.Log("expected:", r, "got:", o.RelocTab[i])
			t.Log(o.RelocTab, o2.RelocTab)
			t.FailNow()
		}
	}
}

func TestRelocateTableMerge(t *testing.T) {
	o1 := vm.NewObject()
	o2 := vm.NewObject()
	o1.AddRelocate(0x0, 0x0000)
	o1.AddRelocate(0x1, 0x0001)
	o2.AddRelocate(0x2, 0x0002)
	o2.AddRelocate(0x3, 0x0003)
	expect := []byte{0x6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1,
		0x2, 0x0, 0x2, 0x3, 0x0, 0x3}

	err := o1.Merge(o2)

	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(o1.RelocTab.Bytes(), expect) {
		t.Fatal("expected:", expect, "got:", o1.RelocTab.Bytes())
	}
}

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

func TestSectionTableMerge(t *testing.T) {
	o1 := vm.NewObject()
	o1.SecTab = vm.SectionTable{
		vm.TEXT: []byte{0x1, 0x2},
		vm.DATA: []byte{0xa, 0xb, 0xc},
	}
	o2 := vm.NewObject()
	o2.SecTab = vm.SectionTable{
		vm.TEXT: []byte{0x3, 0x4},
		vm.DATA: []byte{0xd, 0xe, 0xf},
	}

	if err := o1.Merge(o2); err != nil {
		log.Fatal(err)
	}

	text := []byte{0x1, 0x2, 0x3, 0x4}
	data := []byte{0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

	if !bytes.Equal(o1.SecTab[vm.TEXT], text) {
		t.Fatal("expected", text, "got", o1.SecTab[vm.TEXT])
	}
	if !bytes.Equal(o1.SecTab[vm.DATA], data) {
		t.Fatal("expected", data, "got", o1.SecTab[vm.DATA])
	}
}

func TestSymbol(t *testing.T) {
	o := vm.NewObject()
	o.AddSymbol("main", vm.TEXT, uint16(0x3))
	b := o.SymTab[0].Bytes()
	expect := []byte{0x0, 0x3, byte(vm.TEXT), 0x4, 'm', 'a', 'i', 'n'}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	sym := vm.ScanSymbol(b)
	if sym != o.SymTab[0] {
		t.Log("expected:", o.SymTab[0], "got:", sym)
		t.FailNow()
	}
}

func TestSymbolTable(t *testing.T) {
	o := vm.NewObject()
	o.AddSymbol("foo", vm.DATA, uint16(0xabcd))
	o.AddSymbol("bar", vm.TEXT, uint16(0x1234))

	b := o.SymTab.Bytes()
	expect := []byte{
		0xab, 0xcd, byte(vm.DATA), 0x3, 'f', 'o', 'o',
		0x12, 0x34, byte(vm.TEXT), 0x3, 'b', 'a', 'r',
	}

	if !bytes.Equal(b, expect) {
		t.Log("expected:", expect, "got:", b)
		t.FailNow()
	}

	o2 := vm.NewObject()
	o2.ScanSymbolTable(b)
	for i, s := range o2.SymTab {
		if s != o.SymTab[i] {
			t.Log("expected:", o.SymTab, "got:", o2.SymTab)
			t.FailNow()
		}
	}
}

func TestSymbolTableMerg(t *testing.T) {
	o1 := vm.NewObject()
	o1.AddSymbol("one", vm.TEXT, 0x1)
	o1.AddSymbol("two", vm.TEXT, 0x2)

	o2 := vm.NewObject()
	o2.AddSymbol("three", vm.TEXT, 0x3)
	o2.AddSymbol("four", vm.TEXT, 0x4)

	if err := o1.Merge(o2); err != nil {
		t.Fatal(err)
	}

	if len(o1.SymTab) != 4 {
		t.Fatalf("expected %d symbols, got %d", 4, len(o1.SymTab))
	}

	exp := vm.NewObject()
	exp.AddSymbol("one", vm.TEXT, 0x1)
	exp.AddSymbol("two", vm.TEXT, 0x2)
	exp.AddSymbol("three", vm.TEXT, 0x3)
	exp.AddSymbol("four", vm.TEXT, 0x4)

	for i, sym := range o1.SymTab {
		if sym != exp.SymTab[i] {
			t.Fatal("expected", exp.SymTab[i], "got", sym)
		}
	}

	if err := o1.Merge(o2); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestObject(t *testing.T) {
	o := vm.NewObject()
	o.Entry = 0x03
	o.SecTab = vm.SectionTable{
		vm.TEXT: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa},
		vm.DATA: make([]byte, 0, 0),
	}
	o.AddSymbol("fn", vm.TEXT, 0x0)
	o.AddSymbol("main", vm.TEXT, 0x3)
	o.AddRelocate(0x2, 0x6)

	b := o.Bytes()
	expect := []byte{
		0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf, // magic #
		0x0, 0x3, // entry pt
		0x0, 0x16, // reladdr
		0x0, 0x3, // relsize
		0x0, 0x19, // symaddr
		0x0, 0xe, // symsize
		0x0, 0x27, // secaddr
		0x0, 0xf, // secsize
		0x2, 0x0, 0x6, // reloc1
		0x0, 0x0, byte(vm.TEXT), 0x2, 'f', 'n', // symbol 1
		0x0, 0x3, byte(vm.TEXT), 0x4, 'm', 'a', 'i', 'n', // symbol 2
		0x1,                     // 1 section
		0x0, 0x0, 0x6, 0x0, 0xa, // section text, len 11
		0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, // text
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
	for i := range ob.SecTab {
		if !bytes.Equal(ob.SecTab[i], o.SecTab[i]) {
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

func TestObjectMergeFull(t *testing.T) {
	o1 := vm.NewObject()
	o1.SecTab[vm.TEXT] = []byte{0x0, 0x0}
	o1.SecTab[vm.DATA] = []byte{0x1}
	o1.AddRelocate(0, 0x0)
	o1.AddSymbol("foo", vm.DATA, 0x0)
	o1.AddSymbol("blah", vm.TEXT, 0x0)

	o2 := vm.NewObject()
	o2.Entry = 0x2 // this will fail because merge won't take relocate into account
	o2.SecTab[vm.TEXT] = []byte{0x0, 0x0, 0xff}
	o2.SecTab[vm.DATA] = []byte{0x2}
	o2.AddRelocate(0, 0x0)
	o2.AddSymbol("bar", vm.TEXT, 0x0)

	exp := vm.NewObject()
	exp.Entry = 0x4
	exp.SecTab[vm.TEXT] = []byte{0x0, 0x0, 0x0, 0x2, 0xff}
	exp.SecTab[vm.DATA] = []byte{0x1, 0x2}
	exp.AddRelocate(0, 0x0)
	exp.AddRelocate(3, 0x2)
	exp.AddSymbol("foo", vm.DATA, 0x0)
	exp.AddSymbol("blah", vm.TEXT, 0x0)
	exp.AddSymbol("bar", vm.TEXT, 0x2)

	if err := o1.Merge(o2); err != nil {
		log.Fatal(err)
	}

	if o1.Entry != exp.Entry {
		t.Log("expected entry:", exp.Entry, ", got:", o1.Entry)
		t.Fail()
	}

	// sections
	if !bytes.Equal(o1.SecTab[vm.TEXT], exp.SecTab[vm.TEXT]) {
		t.Log("expected text section:", exp.SecTab[vm.TEXT],
			", got:", o1.SecTab[vm.TEXT])
		t.Fail()
	}
	if !bytes.Equal(o1.SecTab[vm.DATA], exp.SecTab[vm.DATA]) {
		t.Log("expected data section", exp.SecTab[vm.DATA],
			", got:", o1.SecTab[vm.DATA])
		t.Fail()
	}
}
