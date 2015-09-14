package vm_test

/*
func TestEncodeMagicOnly(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	e.Encode([]byte{})
	output := b.Bytes()
	for i := range output {
		if output[i] != []byte{0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xe}[i] {
			t.Log(output)
			t.FailNow()
		}
	}
}

func TestLiteral(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode([]byte("123;0xa;0177")); err != nil {
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

	if err := e.Encode([]byte("0xfff")); err == nil {
		t.Fatal("expected error, got none")
	}

	if err := e.Encode([]byte("xf")); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestInstructions(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode([]byte("nop\njnz\nadd\nsw")); err != nil {
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

	if err := e.Encode([]byte("asdf")); err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestRegisters(t *testing.T) {
	b := new(bytes.Buffer)
	e := vm.NewEncoder(b)
	if err := e.Encode([]byte("$a;$c;$sp")); err != nil {
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

	if err := e.Encode([]byte("$asdf")); err == nil {
		t.Fatal("expected error, got none")
	}
}*/
