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
	Equal     // =
	NotEq     // !=

	// Caracteres especiales
	Lbracket // [
	Rbracket // ]
	Lparen   // (
	Rparen   // )
	Comma    // ,
	OpenQM   // ¿
	CloseQM  // ?
	Dot      // .

	// Palabras reservadas
	Function
	Return
	True
	False
	Null
	And
	Or
	CreateObject
	For
	In
)

// Array con las descripciones de los tokens
var tokenDesc = []string{
	"Illegal",
	"Eof",
	"NewLine",
	"Ident",  // foo, bar
	"Number", // comprende tanto enteros como decimales
	"String", // "foo", 'bar',
	"Assign",

	// Operadores aritméticos
	"Plus",    // +
	"PlusEq",  // +=
	"Minus",   // -
	"MinusEq", // -=
	"Mul",     // *
	"MulEq",   // *=
	"Div",     // /
	"DivEq",   // /=
	"Mod",     // %
	"Pow",     // ^

	// Operadores lógicos
	"Not",       // !
	"Less",      // <
	"LessEq",    // <=
	"Greater",   // >
	"GreaterEq", // >=
	"Equal",     // =
	"NotEq",     // !=

	// Caracteres especiales
	"Lbracket", // [
	"Rbracket", // ]
	"Lparen",   // (
	"Rparen",   // )
	"Comma",    // ,
	"OpenQM",   // ¿
	"CloseQM",  // ?
	"Dot",      // .

	// Palabras reservadas
	"Function",
	"Return",
	"True",
	"False",
	"Null",
	"and",
	"or",
	"CreateObject",
	"For",
	"in",
}

var keywords = map[string]TokenType{
	"Func":         Function,
	"Return":       Return,
	"True":         True,
	"False":        False,
	"Null":         Null,
	"and":          And,
	"or":           Or,
	"CreateObject": CreateObject,
	"For":          For,
	"in":           In,
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

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}

func GetTokenStr(t TokenType) string {
	return tokenDesc[t]
}
