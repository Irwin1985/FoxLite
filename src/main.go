package main

import (
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
				repl.RunFile(sourcePath)
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
		repl.RunPrompt(os.Stdin, os.Stdout)
	}
}

func printVersion() string {
	return "v1.0.1"
}
