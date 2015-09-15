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
	buf *bytes.Buffer
	ob  Object
	//stab map[string]uint16
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{Writer: w, buf: new(bytes.Buffer)} //,
	//		stab: make(map[string]uint16)}
}

func (e *Encoder) Encode(src []byte) error {
	f, errs := Parse(src)
	// TODO need error list satisfying error interface
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		log.Fatal("failed to assemble file")
	}

	// generate text & data bytes
	if err := e.file(f); err != nil {
		//log.Fatal(err)
		return err
	}

	_, err := e.Write(e.ob.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (e *Encoder) emit(b ...byte) {
	n, err := e.buf.Write(b)
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
				//e.stab[k] = uint16(e.buf.Len())
				addr := uint16(e.buf.Len())
				e.ob.SymTab = append(e.ob.SymTab,
					Symbol{Name: k, Addr: addr})
				e.ob.RelocTab = append(e.ob.RelocTab,
					Relocate{Offset: addr, SymIndex: byte(len(e.ob.SymTab))})
				fmt.Println("new symbol:", k, addr)
				e.sub(v)
			}
			e.ob.SecTab = append(e.ob.SecTab,
				OSection{Name: "text", Data: e.buf.Bytes()})
		default:
			log.Fatal("unexpected section type")
		}
	}
	//fmt.Println(e.stab)
	return
}

func (e *Encoder) sub(il []*Instruction) {
	for _, i := range il {
		switch i.Op {
		case JMP, JPZ, JNZ, CALL:
			// TODO replace
			s, ok := e.ob.SymTab.Lookup(i.Value) //e.stab[i.Value]
			if !ok {
				log.Fatal("undeclared symbol", i.Value)
			}
			b := toBytes(s.Addr)
			e.emit(byte(i.Op), b[0], b[1])
		case MVI:
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
