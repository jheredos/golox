package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jheredos/golox/lox"
)

var hadError bool = false

// here's a change

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	fmt.Println("runFile", path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// run(string(bytes))

	// tokens := lox.Scan(string(bytes))
	tokens := lox.Lex(string(bytes))

	fmt.Println("Lexing complete")

	for i := 0; i < len(tokens); i++ {
		// fmt.Println(tokens[i].Lexeme)
		fmt.Printf("Line %d: type %d, value \"%s\"\n", tokens[i].Line, tokens[i].Type, tokens[i].Lexeme)
	}

	if hadError {
		os.Exit(1)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		// tokens := lox.Scan(line)
		tokens := lox.Lex(line)

		for i := 0; i < len(tokens); i++ {
			fmt.Println(tokens[i].Lexeme)
		}

		if err != nil {
			fmt.Println(err)
		}

		if err := run(line); err != nil {
			fmt.Println(err)
		}

		hadError = false
	}
}

func run(src string) error {
	fmt.Print(src)

	return nil
}

func newError(line int, message string) {
	reportError(line, "", message)
}

func reportError(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true
}
