package token

import "fmt"

type TokenType int

const (
	Illegal = iota
	Eof
	NewLine
	Ident  // foo, bar
	Number // comprende tanto enteros como decimales
	String // "foo", 'bar',
	Assign

	// Operadores aritméticos
	Plus    // +
	PlusEq  // +=
	Minus   // -
	MinusEq // -=
	Mul     // *
	MulEq   // *=
	Div     // /
	DivEq   // /=
	Mod     // %
	Pow     // ^

	// Operadores lógicos
	Not       // !
	Less      // <
	LessEq    // <=
	Greater   // >
	GreaterEq // >=
)

// Array con las descripciones de los tokens
var tokenDesc = []string{
	"Illegal",
	"Eof",
	"NewLine",
	"Ident",
	"Number",
	"String",
	"Assign",

	// Operadores aritméticos
	"Plus",
	"PlusEq",
	"Minus",
	"MinusEq",
	"Mul",
	"MulEq",
	"Div",
	"DivEq",
	"Mod",
	"Pow",

	// Operadores lógicos
	"Not",
	"Less",
	"LessEq",
	"Greater",
	"GreaterEq",
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

func (t *Token) Str() string {
	return fmt.Sprintf("<%s, '%s'> at [%d:%d]", tokenDesc[t.Type], t.Literal, t.Line, t.Col)
}
