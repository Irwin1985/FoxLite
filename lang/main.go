package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run fibonacci.prg
	input := `IIF(1 == 1, "POLICIA 1", "POLICIA 2")`
	mode := "repl"
	repl.Start(mode, input)
}
