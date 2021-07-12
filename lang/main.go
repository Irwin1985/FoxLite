package main

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/interpreter"
	"FoxLite/lang/lexer"
	"FoxLite/lang/parser"
	"FoxLite/lang/token"
	"fmt"
)

func main() {
	//testLexer()
	//testStream()
	//testParser()
	testInterpreter()
}

func testInterpreter() {
	input := `
		IF .T. THEN
			RETURN "ES TRUE"
		ELSE
			RETURN "ES FALSE"
		ENDIF
	`
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	if len(p.Errors()) > 0 {
		printErrors(p.Errors())
		return
	}
	i := interpreter.NewInterpreter(program)
	out := i.Interpret()
	fmt.Printf("%v", out)
}

func testParser() {
	input := `
		FUNC ADD(X, Y)
			LOCAL A = 10
			IF A >= 20 THEN
				A = 20
			ELSE
				A = 30
			ENDIF
			IF A <= 20
				A = 30
			ENDIF
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
