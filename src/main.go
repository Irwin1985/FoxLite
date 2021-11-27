package main

import (
	"FoxLite/src/lexer"
	"FoxLite/src/repl"
	"fmt"
	"os"
	"path/filepath"
)

const FILENAME = 2

func main() {
	if len(os.Args) >= 2 {
		cmd := os.Args[1]
		switch cmd {
		case "run":
			fileName := os.Args[FILENAME]
			if sourcePath, err := filepath.Abs(fileName); err == nil {
				runFile(sourcePath)
			} else {
				fmt.Printf("No such file: '%s'", fileName)
			}
		case "compile":
			fmt.Println("policia compile")
		case "fmt":
			fmt.Println("policia format")
		default:
			fmt.Printf("unknown command %s\n", cmd)
			os.Exit(1)
		}
	} else {
		runPrompt()
	}
}

func runFile(fileName string) {
	l := lexer.New()
	l.ScanFile(fileName)
	repl.Execute(l)
}

func runPrompt() {
	repl.Start(os.Stdin, os.Stdout)
}
