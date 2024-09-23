package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pmqueiroz/umbra/interpreter"
)

func header() {
	fmt.Print("Welcome to Umbra REPL!\nEnter :q to exit.\n")
}

func Repl(evaluate func(content string, env *interpreter.Environment)) {
	header()
	reader := bufio.NewReader(os.Stdin)
	env := interpreter.NewEnvironment(nil)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil || line == ":q\n" {
			break
		}

		evaluate(strings.Join([]string{"module repl", line}, "\n"), env)
	}
}
