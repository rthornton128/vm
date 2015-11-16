package vm_test

/*
func TestEncodeMinimal(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode(strings.NewReader("")); err != nil {
		log.Fatal(err)
	}

	expect := []byte{
		0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf, // magic
		0x0, 0x0, // entry
		0x0, 0x0, // relocation table
		0x0, 0x0, // symbol table
		0x0, 0x0, // sections
	}
	if !bytes.Equal(b.Bytes(), expect) {
		t.Fatal("expected", expect, "got", b.Bytes())
	}
}

func TestLiteral(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode(strings.NewReader("123;0xa;0177")); err != nil {
		t.Fatal(err)
	}
	output := b.Bytes()
	expected := []byte{0x7b, 0xa, 0x7f}
	for i := range output[8:] {
		if output[i+8] != expected[i] {
			t.Log(output[8:])
			t.Log("expected:", expected[i], "got", output[i+8])
			t.FailNow()
		}
	}

	if err := e.Encode(strings.NewReader("0xfff")); err == nil {
		t.Fatal("expected error, got none")
	}

	if err := e.Encode(strings.NewReader("xf")); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestInstructions(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode(strings.NewReader("nop\njnz\nadd\nsw")); err != nil {
		t.Fatal(err)
	}
	output := b.Bytes()
	expected := []vm.Instruction{vm.NOP, vm.JNZ, vm.ADD, vm.SW}
	for i := range output[8:] {
		if vm.Instruction(output[i+8]) != expected[i] {
			t.Log(output[8:])
			t.Log("expected:", expected[i], "got", output[i+8])
			t.FailNow()
		}
	}

	if err := e.Encode(strings.NewReader("asdf")); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestRegisters(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode(strings.NewReader("$a;$c;$sp")); err != nil {
		t.Fatal(err)
	}
	output := b.Bytes()
	expected := []vm.Register{vm.REGA, vm.REGC, vm.SP}
	for i := range output[8:] {
		if vm.Register(output[i+8]) != expected[i] {
			t.Log(output[8:])
			t.Log("expected:", expected[i], "got", output[i+8])
			t.FailNow()
		}
	}

	if err := e.Encode(strings.NewReader("$asdf")); err == nil {
		t.Fatal("expected error, got none")
	}
}*/
