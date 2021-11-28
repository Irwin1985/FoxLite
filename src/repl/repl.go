package repl

import (
	"FoxLite/src/evaluator"
	"FoxLite/src/lexer"
	"FoxLite/src/object"
	"FoxLite/src/parser"
	"FoxLite/src/token"
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

const PROMPT = ">>> "
const VERSION = "1.0.1"

func RunFile(fileName string) { // Ejecuta el código de un fichero.
	env := createEnvironment()
	l := lexer.New()
	l.ScanFile(fileName)
	Execute(l, os.Stdout, env)
}

func RunPrompt(in io.Reader, out io.Writer) {
	displayWelcome()
	scanner := bufio.NewScanner(in)
	env := createEnvironment()

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
		Execute(l, out, env)
	}
}

func Execute(l *lexer.Lexer, out io.Writer, env *object.Environment) {
	dumpTokens := false
	if dumpTokens {
		tok := l.NextToken()
		for tok.Type != token.Eof {
			fmt.Println(tok.Str())
			tok = l.NextToken()
		}
		fmt.Println(tok.Str())
	} else {
		p := parser.New(l)
		program := p.Parse()
		errors := p.Errors()
		if len(errors) > 0 {
			printErrors(errors, out)
			return
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if evaluated.Type() == object.ErrorObj {
				msg := fmt.Sprintf("%s %s\n", l.GetErrorFormat(nil), evaluated.Inspect())
				fmt.Println(msg)
				return
			}
			fmt.Println(evaluated.Inspect())
		}
	}
	//fmt.Println(color.Green + "policia execute 2!" + color.Reset)
}

func printErrors(errors []string, out io.Writer) {
	for _, msg := range errors {
		_, err := io.WriteString(out, msg)
		if err != nil {
			panic(err)
		}
	}
}

func displayWelcome() {
	fmt.Printf("Welcome to Foxlite programming language, version %s\n", VERSION)
	fmt.Printf("Data and Time %s\n", time.Now().Format(time.Stamp))
	fmt.Printf("Type \"quit\" to exit or \"help\" for more information.\n")
}

func createEnvironment() *object.Environment {
	e := object.NewEnv()
	return e
}
