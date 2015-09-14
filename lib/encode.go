package vm

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
)

type symbol struct {
	index uint16
	size  uint16
}

type Encoder struct {
	io.Writer
	buf  *bytes.Buffer
	stab map[string]uint16
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Writer: w, buf: new(bytes.Buffer),
		stab: make(map[string]uint16)}
}

func addrOf(i int16) []byte {
	return []byte{byte(i >> 8), byte(i & 0xff)}
}

func (e *Encoder) Encode(src []byte) error {
	f, errs := Parse(src)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		log.Fatal("failed to assemble file")
	}

	// generate text & data bytes
	if err := e.file(f); err != nil {
		log.Fatal(err)
	}

	// magic number
	e.Write(e.Writer, 0xd, 0xe, 0xa, 0xd, 0xb, 0xe, 0xe, 0xf)

	// entry point
	msb, lsb := toBytes(e.stab["main"])
	e.Write(e.Writer, msb, lsb)

	// symbol table offset
	//e.Write(addrOf(e.sz + 12))

	// program text
	e.Write(e.Writer, e.buf.Bytes()...)

	// symbol table
	//for _, v := range e.stab {
	//e.Write(addrOf(v)) // index
	//e.Write(addrOf(0)) // address
	//}

	return nil
}

/*
func (e *Encoder) data(d *secData) {
	if d != nil {
		for _, v := range d.data {
			e.emit([]byte{v})
		}
	}
}
*/
func (e *Encoder) emit(b ...byte) {
	e.Write(e.buf, b...)
}

func (e *Encoder) Write(w io.Writer, b ...byte) {
	n, err := w.Write(b)
	if n != len(b) {
		log.Println("failed to write all bytes")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Encoder) file(f *File) error {
	for k, s := range f.sections[0].(*TextSection).m {
		e.stab[k] = uint16(e.buf.Len())
		if k == "main" {
			defer e.sub(s)
			continue
		}
		e.sub(s)
	}
	fmt.Println(e.stab)
	return nil
}

func (e *Encoder) sub(il []*Instruction) {
	for _, i := range il {
		switch i.Op {
		case JMP, CALL:
			index, ok := e.stab[i.Value]
			if !ok {
				log.Fatal("undefined symbol", i.Value)
			}
			msb, lsb := toBytes(index)
			e.emit(byte(i.Op), msb, lsb)
		case JPZ, JNZ, MVI:
			v, err := strconv.ParseInt(i.Value, 0, 8)
			if err != nil {
				log.Fatal(err)
			}
			e.emit(byte(i.Op), byte(v))
		default:
			e.emit(byte(i.Op))
		}
	}
	return
}
