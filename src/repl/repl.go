package repl

import (
	"FoxLite/src/lexer"
	"FoxLite/src/token"
	"bufio"
	"fmt"
	"io"
	"time"
)

const PROMPT = ">>> "
const VERSION = "1.0.1"

func Start(in io.Reader, _ io.Writer) {
	displayWelcome()
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		input := scanner.Text()
		if len(input) == 0 {
			continue
		}
		if input == "quit" {
			break
		}
		l := lexer.New()
		l.ScanText([]rune(input))
		Execute(l)
	}
}

func displayWelcome() {
	fmt.Printf("Welcome to Foxlite programming language, version %s\n", VERSION)
	fmt.Printf("Data and Time %s\n", time.Now().Format(time.Stamp))
	fmt.Printf("Type \"quit\" to exit or \"help\" for more information.\n")
}

func Execute(l *lexer.Lexer) {
	tok := l.NextToken()
	for tok.Type != token.Eof {
		fmt.Println(tok.Str())
		tok = l.NextToken()
	}
	fmt.Println(tok.Str())
	//fmt.Println(color.Green + "policia execute 2!" + color.Reset)
}
