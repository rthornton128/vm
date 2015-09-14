package vm

type Section interface {
	//Name() string
	//Data() interface{}
}

type (
	File struct {
		sections []Section
		//magic []byte
		stab map[string]int
	}
	Instruction struct {
		Op    Opcode
		Value string
	}
	DataSection struct {
		m map[string]byte
	}
	TextSection struct {
		m map[string][]*Instruction
	}
)
