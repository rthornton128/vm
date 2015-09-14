package vm

import "github.com/rthornton128/gompto/lex"

var symbols = map[string]lex.Token{
	"%": PERCENT,
	".": DOT,
	":": COLON,
	"$": DOLLAR,
}
