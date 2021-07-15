package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run fibonacci.prg
	input := `hola`
	mode := "repl"
	repl.Start(mode, input)
}
