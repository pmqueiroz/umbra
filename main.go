package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/tokens"
)

func run(contents ...string) {
	content := strings.Join(contents, "\n")
	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	module := ast.Parse(tokens)

	fmt.Printf("%#v\n", module)
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

		run(
			"module main",
			line,
		)
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
