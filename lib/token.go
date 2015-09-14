package vm

import "github.com/rthornton128/gompto/lex"

const (
	PERCENT lex.Token = iota + lex.TokenStart
	DOT
	COLON
	DOLLAR
)
