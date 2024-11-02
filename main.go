package main

import (
	"fmt"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/cli"
	"github.com/pmqueiroz/umbra/helpers"
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
		fmt.Println(err)
	}
}

func main() {
	args := cli.Parse()

	if args.Input != "" {
		content, err := helpers.ReadFile(args.Input)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		env := interpreter.NewEnvironment(nil)

		run(content, RunOptions{
			Options: args.Options,
			Env:     env,
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
