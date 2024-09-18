package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/tokens"
)

func run(content string) {
	tokens, err := tokens.Tokenizer(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	ast.Parse(tokens)
}

func runFile(path string) {
	fileContent, err := readFile(path)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	run(fileContent)
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Welcome to Umbra REPL!\nEnter :q to exit.\n")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil || line == ":q\n" {
			break
		}

		run(line)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: umbra [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
