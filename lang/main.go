package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run fibonacci.prg
	input := `
		LOCAL A = 20
		A += 10
	`
	mode := "repl"
	repl.Start(mode, input)
}
