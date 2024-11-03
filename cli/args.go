package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	PrintAst    bool
	ShowVersion bool
}

type Args struct {
	Options Options
	Path    string
}

const HELP_HEADER = `Usage: umbra [options] [[file] [arguments]]

Options:`

func Parse() Args {
	parsedArgs := Args{}

	flag.BoolVar(&parsedArgs.Options.PrintAst, "ast", false, "Prints the AST of the program")
	flag.BoolVar(&parsedArgs.Options.ShowVersion, "version", false, "Display the version of umbra")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, HELP_HEADER)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	args = filterArgs(args)

	if len(args) > 0 {
		parsedArgs.Path = args[0]
	}

	return parsedArgs
}

func filterArgs(input []string) []string {
	var result []string
	for _, str := range input {
		if !strings.HasPrefix(str, "-") {
			result = append(result, str)
		}
	}
	return result
}
