package vm

import (
	"errors"
	"fmt"
	"log"

	"github.com/rthornton128/gompto/lex"
)

type Parser struct {
	lexer  *lex.Lex
	errors []error
	offset int

	lit string
	pos int
	tok lex.Token
}

func Parse(src []byte) (*File, []error) {
	p := newParser(string(src))
	f := p.parseFile()
	return f, p.errors
}

func newParser(src string) *Parser {
	s := new(lex.Scanner)
	s.Init(src)
	l := lex.NewLex(s)
	l.Symbols = symbols
	p := &Parser{
		lexer:  l,
		errors: make([]error, 0),
	}
	p.next()
	return p
}

func (p *Parser) error(args ...interface{}) {
	p.errors = append(p.errors,
		errors.New(fmt.Sprint("line: ", p.pos, " - ", args)))
}

func (p *Parser) expect(t lex.Token) {
	if p.tok != t {
		//fmt.Println("expected", t, "got:", p.tok, " (", p.lit, ")")
		p.error("expected", t, "got:", p.tok, " (", p.lit, ")")
	}
	p.next()
}

func (p *Parser) ident() string {
	l := p.lit
	p.expect(lex.IDENT)
	return l
}

func (p *Parser) instruction(id string) *Instruction {
	i, err := LookupOpcode(id)
	if err != nil {
		log.Fatal(err)
	}

	switch i {
	case JMP, CALL:
		p.expect(DOLLAR)
		return &Instruction{Op: i, Value: p.ident()}
	case JPZ, JNZ, MVI:
		return &Instruction{Op: i, Value: p.literal()}
	case CLA, INC, NOP, POP, PUSH, RET:
		return &Instruction{Op: i}
	default:
		return &Instruction{Op: i | Opcode(p.register())}
	}
}

func (p *Parser) literal() string {
	l := p.lit
	p.expect(lex.INT)
	return l
}

func (p *Parser) next() {
	p.lit, p.tok, p.pos = p.lexer.Lex()
	fmt.Println("next:", p.lit, p.tok, p.pos)
}

func (p *Parser) register() Register {
	p.expect(PERCENT)
	r, err := LookupRegister(p.lit)
	if err != nil {
		p.error(err)
	}
	p.next()
	return r
}

func (p *Parser) parseFile() *File {
	sections := make([]Section, 0)
	for p.tok != lex.EOF {
		p.expect(DOT)
		ident := p.ident()
		switch ident {
		case "text":
			sections = append(sections, p.sectionText())
		default:
			p.error("expected valid section name, got", p.lit)
			return nil
		}
	}

	return &File{sections: sections}
}

/*
func (p *Parser) sectionData() *DataSection {
	// parse label and literal pairs until next section marker found
	data := make(map[string]byte)
	for p.tok == IDENT { //&& p.tok != EOF {
		lab := p.ident()
		p.expect(COLON)
		data[lab] = p.literal()
		//p.stab[lab] = p.offset
		//p.offset++
	}
	return &DataSection{m: data}
}*/

func (p *Parser) sectionText() *TextSection {
	text := make(map[string][]*Instruction)
	var sub string
	for p.tok == lex.IDENT {
		id := p.ident()
		if p.tok == COLON { // new subroutine
			sub = id
			p.next()
			continue
		}
		text[sub] = append(text[sub], p.instruction(id))
	}
	return &TextSection{m: text}
}
