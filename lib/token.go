package vm

import "github.com/rthornton128/gct/lex"

const (
	PERCENT lex.Token = iota + lex.TokenStart
	DOT
	COLON
	DOLLAR
)
