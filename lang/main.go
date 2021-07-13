package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run F:\Desarrollo\GitHub\GOPATH\src\FoxLite\lang\samples\fibonacci.prg
	input := `
		LOCAL(
			A = 10,
			B = 20,
			C = A + B
		)
		RETURN C
	`
	mode := "interpreter"
	repl.Start(mode, input)
}
