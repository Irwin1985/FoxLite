package token

import (
	"fmt"
	"strings"
)

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
	If
	Else
	Do
	Case
	Otherwise
	EndCase
	While
	Exit
	Loop
	// Variables
	Private // Private
	Local   // Local
	Public  // Public
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
	"If",
	"Else",
	"Do",
	"Case",
	"Otherwise",
	"EndCase",
	"While",
	"Exit",
	"Loop",
	"Private",
	"Local",
	"Public",
}

var keywords = map[string]TokenType{
	"func":         Function,
	"function":     Function,
	"return":       Return,
	"true":         True,
	"false":        False,
	"null":         Null,
	"and":          And,
	"or":           Or,
	"createobject": CreateObject,
	"for":          For,
	"in":           In,
	"if":           If,
	"else":         Else,
	"do":           Do,
	"case":         Case,
	"otherwise":    Otherwise,
	"endcase":      EndCase,
	"while":        While,
	"exit":         Exit,
	"loop":         Loop,
	"prv":          Private,
	"loc":          Local,
	"pub":          Public,
	"private":      Private,
	"local":        Local,
	"public":       Public,
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
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return Ident
}

func GetTokenStr(t TokenType) string {
	return tokenDesc[t]
}
