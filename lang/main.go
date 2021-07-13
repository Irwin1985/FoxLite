package main

import (
	"FoxLite/lang/repl"
)

func main() {
	input := `run F:\Desarrollo\GitHub\GOPATH\src\FoxLite\lang\samples\fibonacci.prg`
	mode := "repl"
	repl.Start(mode, input)
}
