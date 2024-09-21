package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/cli"
	"github.com/pmqueiroz/umbra/interpreter"
	"github.com/pmqueiroz/umbra/tokens"
)

func run(options cli.Options, contents ...string) {
	content := strings.Join(contents, "\n")
	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	module := ast.Parse(tokens)

	if options.PrintAst {
		cli.PrintAst(module)
	}

	env := interpreter.NewEnvironment(nil)

	if err := interpreter.Interpret(module, env); err != nil {
		fmt.Println("Erro:", err)
	}
}

func runFile(path string, options cli.Options) {
	fileContent, err := readFile(path)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	run(options, fileContent)
}

func runPrompt(options cli.Options) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Welcome to Umbra REPL!\nEnter :q to exit.\n")

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil || line == ":q\n" {
			break
		}

		run(
			options,
			"module main",
			line,
		)
	}
}

func main() {
	args := cli.Parse()

	if args.Input != "" {
		runFile(args.Input, args.Options)
	} else {
		runPrompt(args.Options)
	}
}
