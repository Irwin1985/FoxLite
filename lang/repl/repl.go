package repl

import (
	"FoxLite/lang/ast"
	"FoxLite/lang/interpreter"
	"FoxLite/lang/lexer"
	"FoxLite/lang/object"
	"FoxLite/lang/parser"
	"FoxLite/lang/token"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const PROMPT = `
 /\   /\   
//\\_//\\     ____
\_     _/    /   /
 / * * \    /^^^]
 \_\O/_/    [   ]
  /   \_    [   /
  \     \_  /  /
   [ [ /  \/ _/
  _[ [ \  /_/
`

const ERROR = `
 ^...^
<_@ @_>   
  \_/
`

const VERSION = "1.0"

var globalEnv = object.NewEnvironment()

func Start(mode string, input string) {
	if mode == "repl" {
		repl()
	} else if mode == "lexer" {
		debugLexer(input)
	} else if mode == "parser" {
		debugParser(input)
	} else if mode == "interpreter" {
		debugInterpreter(input)
	}
}

func repl() {
	displayWelcome()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		input := scanner.Text()
		if len(input) < 0 {
			continue
		}
		if input == "quit" {
			break
		}
		evalInput(input)
	}
}

func evalInput(input string) {
	if input[0:3] == "run" {
		filePath := strings.TrimSpace(input[3:])
		err := runFile(filePath)
		if err != nil {
			fmt.Printf("%s\n%v\n", ERROR, err)
		}
	} else {
		run(input)
	}
}

func run(input string) error {
	// measuring time
	start := time.Now()
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	errors := p.Errors()
	if len(errors) > 0 {
		printErrors(errors)
	}
	i := interpreter.NewInterpreter(program, globalEnv)
	output := i.Interpret()
	elapsed := time.Since(start)
	timeStr := fmt.Sprintf("Elapsed time: %s\n", elapsed)
	if output != nil {
		switch obj := output.(type) {
		case *object.Error:
			fmt.Printf("%s\n%s\n", ERROR, obj.Message)
		default:
			fmt.Println(obj)
		}
		// printing the elapsed time
		fmt.Print(timeStr)
	}
	return nil
}

func runFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("No souch file: %s\n", filePath)
	}
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return run(string(fileContent))
}

func debugLexer(input string) {
	l := lexer.NewLexer(input)
	tok := l.NextToken()
	for tok.Type != token.EOF {
		fmt.Println(tok.ToString())
		tok = l.NextToken()
	}
}

func debugParser(input string) {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.Parse()
	errors := p.Errors()
	if len(errors) > 0 {
		printErrors(errors)
	}
	print := ast.NewAstPrinter(program)
	fmt.Printf("%v\n", print.PrettyPrint())
}

func debugInterpreter(input string) {
	evalInput(input)
}

/**********************************************************
* HELPER FUNCTIONS
***********************************************************/
func displayWelcome() {
	fmt.Printf("%s\n", PROMPT)
	fmt.Printf("Welcome to FoxLite Version: %s\n", VERSION)
	fmt.Printf("Data and time %v\n", time.Now().Format(time.ANSIC))
	fmt.Printf("Type 'quit' to exit.\n")
}

func printErrors(errors []string) {
	fmt.Printf("%s\n", ERROR)
	for _, msg := range errors {
		fmt.Printf("%s\n", msg)
	}
}
