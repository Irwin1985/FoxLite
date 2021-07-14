package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run fibonacci.prg
	input := `
		DO CASE
		CASE 1 > 1
			RETURN "POLICIA 1"
		CASE 1 < 1
			RETURN "POLICIA 2"
		CASE 1 != 1
			RETURN "POLICIA 3"
		ENDCASE
	`
	mode := "interpreter"
	repl.Start(mode, input)
}
