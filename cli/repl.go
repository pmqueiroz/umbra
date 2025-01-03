package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/environment"
)

func header() {
	fmt.Print("Welcome to Umbra REPL!\nEnter :q to exit.\n")
}

func Repl(evaluate func(content string, env *environment.Environment)) {
	header()
	reader := bufio.NewReader(os.Stdin)
	env := environment.NewEnvironment(nil)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil || line == ":q\n" {
			break
		}

		evaluate(line, env)
	}
}
