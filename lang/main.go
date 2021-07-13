package main

import (
	"FoxLite/lang/repl"
)

func main() {
	//run F:\Desarrollo\GitHub\GOPATH\src\FoxLite\lang\samples\fibonacci.prg
	input := `
		FUNCTION FIRST_NAME(FNAME)
			FUNCTION LAST_NAME(LNAME)
				RETURN FNAME + ", " + LNAME
			ENDFUNC
			RETURN LAST_NAME
		ENDFUNC		
		RETURN LAST_NAME("RODRIGUEZ")
	`
	mode := "interpreter"
	repl.Start(mode, input)
}
