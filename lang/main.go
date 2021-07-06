package main

import (
	"FoxLite/lang/lexer"
	"FoxLite/lang/parser"
	"FoxLite/lang/token"
	"fmt"
)

func main() {
	//testLexer()
	//testStream()
	testParser()
}

func testParser() {
	input := `	
	IF .T.
		.t.
	ENDIF
`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Program()
	fmt.Println(program.String())
}

func testLexer() {
	input := `
	LOCAL A = 10
	PRIVATE B = 20
	PUBLIC C = A + B
	IF A > B THEN
		MESSAGEBOX("a es mayor")
	ELSE
		MESSAGEBOX("b es mayor")
	ENDIF
	
	* Sample function
	FUNCTION DOUBLE(X, Y)
		RETURN X + Y
	ENDFUNC

	FUNCTION TRIPE(X, DOUBLE)
		RETURN DOUBLE(X) + X
	ENDFUNC
`
	l := lexer.NewLexer(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		fmt.Println(tok.ToString())
		tok = l.NextToken()
	}
	fmt.Println(tok.ToString())
}

func testStream() {
	input := `if a > '0' then a else b;`
	s := lexer.NewStream(input)
	c := s.Read()
	for c != lexer.EOF_CHAR {
		fmt.Printf("%c", c)
		c = s.Read()
	}
}
