package vm

import "github.com/rthornton128/gct/lex"

var symbols = map[string]lex.Token{
	"%": PERCENT,
	".": DOT,
	":": COLON,
	"$": DOLLAR,
}
