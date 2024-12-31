package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/cli"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/interpreter"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
)

// set through ldflags
var Version = "development"

type RunOptions struct {
	cli.Options
	Env *environment.Environment
}

func run(content string, options RunOptions) error {
	tokens, err := tokens.Tokenize(content)

	if err != nil {
		return err
	}

	if options.PrintTokens {
		cli.PrintTokens(tokens)
	}

	module := ast.Parse(tokens)

	if options.PrintAst {
		cli.PrintAst(module)
	}

	if err := interpreter.Interpret(module, options.Env); err != nil {
		return err
	}

	return nil
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
			fmt.Println(exception.NewUmbraError("GN001", nil, args.Path))
			return
		}
		content, err := helpers.ReadFile(args.Path)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		env := environment.NewEnvironment(nil)

		env.Create(nil, "__FILE__", __FILE__, types.STR, false, false, false)

		runErr := run(content, RunOptions{
			Options: args.Options,
			Env:     env,
		})

		if runErr != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		cli.Repl(func(content string, env *environment.Environment) {
			runErr := run(content, RunOptions{
				Options: args.Options,
				Env:     env,
			})

			if runErr != nil {
				fmt.Println(runErr)
			}
		})
	}
}
