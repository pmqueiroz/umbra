package main

import (
	"fmt"
	"path/filepath"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/cli"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/interpreter"
	"github.com/pmqueiroz/umbra/tokens"
)

// set through ldflags
var Version = "development"

type RunOptions struct {
	cli.Options
	Env *environment.Environment
}

func run(content string, options RunOptions) {
	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	if options.PrintTokens {
		cli.PrintTokens(tokens)
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

	if args.Options.ShowVersion {
		fmt.Printf("%s\n", Version)
		return
	}

	if args.Path != "" {
		__FILE__, err := filepath.Abs(args.Path)
		if err != nil {
			fmt.Println(exception.NewGenericError("GN001", args.Path))
			return
		}
		content, err := helpers.ReadFile(args.Path)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		env := environment.NewEnvironment(nil)

		env.Create("__FILE__", __FILE__, tokens.STR_TYPE, false, false)

		run(content, RunOptions{
			Options: args.Options,
			Env:     env,
		})
	} else {
		cli.Repl(func(content string, env *environment.Environment) {
			run(content, RunOptions{
				Options: args.Options,
				Env:     env,
			})
		})
	}
}
