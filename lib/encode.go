package vm

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Encoder struct {
	io.Writer
	buf  *bytes.Buffer
	ob   Object
	stab map[string]uint16
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Writer: w, buf: new(bytes.Buffer),
		stab: make(map[string]uint16)}
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

	_, err := e.Write(e.ob.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (e *Encoder) emit(b ...byte) {
	n, err := e.Write(b)
	if n != len(b) {
		log.Println("failed to write all bytes")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (e *Encoder) file(f *File) error {
	if len(f.sections) > 0 {
		e.sections(f.sections)
	}
	return nil
}

// TODO horrific, section handling needs massive (re)work
func (e *Encoder) sections(secs []Section) {
	for _, s := range secs {
		switch x := s.(type) {
		case *TextSection:
			for k, v := range x.m {
				e.stab[k] = uint16(e.buf.Len())
				if k == "main" {
					defer e.sub(v)
					continue
				}
				e.sub(v)
			}
		default:
			log.Fatal("unexpected section type")
		}
	}
	fmt.Println(e.stab)
	return
}

func (e *Encoder) sub(il []*Instruction) {
	for _, i := range il {
		switch i.Op {
		case JMP, CALL:
			index, ok := e.stab[i.Value]
			if !ok {
				log.Fatal("undefined symbol", i.Value)
			}
			b := toBytes(index)
			e.emit(byte(i.Op), b[0], b[1])
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
