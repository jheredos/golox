package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jheredos/golox/lox"
)

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

	tokens := lox.Lex(string(bytes))
	_, err = lox.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// run(ast)
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		tokens := lox.Lex(line)

		_, err = lox.Parse(tokens)
		if err != nil {
			fmt.Println(err)
		}

		if err := run(line); err != nil {
			fmt.Println(err)
		}
	}
}

func run(src string) error {
	fmt.Print(src)

	return nil
}
