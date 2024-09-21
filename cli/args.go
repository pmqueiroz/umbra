package cli

import (
	"flag"
	"fmt"
	"os"
)

type Options struct {
	PrintAst bool
}

type Args struct {
	Options   Options
	Input     string
	Arguments []string
}

const HELP_HEADER = `Usage: umbra [options] [[file] [arguments]]

Options:`

func Parse() Args {
	parsedArgs := Args{}

	flag.BoolVar(&parsedArgs.Options.PrintAst, "ast", false, "Prints the AST of the program")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, HELP_HEADER)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		parsedArgs.Input = args[0]
		parsedArgs.Arguments = args[1:]
	}

	return parsedArgs
}
