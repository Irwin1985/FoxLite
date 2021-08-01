package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run fibonacci.prg
	input := `		
		LOCAL PI = 3.14
	`
	mode := "lexer"
	repl.Start(mode, input)
}
