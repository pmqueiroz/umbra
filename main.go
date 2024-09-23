package main

import (
	"fmt"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/cli"
	"github.com/pmqueiroz/umbra/interpreter"
	"github.com/pmqueiroz/umbra/tokens"
)

type RunOptions struct {
	cli.Options
	Env *interpreter.Environment
}

func run(content string, options RunOptions) {
	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	module := ast.Parse(tokens)

	if options.PrintAst {
		cli.PrintAst(module)
	}

	if err := interpreter.Interpret(module, options.Env); err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	args := cli.Parse()

	if args.Input != "" {
		content, err := readFile(args.Input)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		run(content, RunOptions{
			Options: args.Options,
			Env:     interpreter.NewEnvironment(nil),
		})
	} else {
		cli.Repl(func(content string, env *interpreter.Environment) {
			run(content, RunOptions{
				Options: args.Options,
				Env:     env,
			})
		})
	}
}
