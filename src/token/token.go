package token

type TokenType string

const (
	Illegal = iota
	Eof
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

	// Operadores lógicos
	Not       // !
	Less      // <
	LessEq    // <=
	Greater   // >
	GreaterEq // >=

)
