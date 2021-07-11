package main

import (
	"FoxLite/lang/ast"
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
		FUNC ADD(X, Y)
			LOCAL A = 10
			RETURN A
		ENDFUNC
	`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
		return
	}
	a := ast.NewAstPrinter(program)
	out := a.PrettyPrint()
	fmt.Println(out)
}

func testLexer() {
	input := `
	LOCAL A = 10
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

func printErrors(errors []string) {
	fmt.Println("Error:")
	for _, msg := range errors {
		fmt.Println(msg)
	}
}
