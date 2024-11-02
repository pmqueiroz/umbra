package interpreter

import (
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/tokens"
)

func ResolveModule(module string) (string, error) {
	content, err := helpers.ReadFile(fmt.Sprintf("%s/lib/%s.u", os.Getenv("UMBRA_PATH"), module))

	return content, err
}

func LoadModule(path string) (Environment, error) {
	content, err := ResolveModule(path)

	if err != nil {
		return Environment{}, fmt.Errorf("unable to include %s. module does not exits", path)
	}

	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	module := ast.Parse(tokens)

	namespace := NewEnvironment(nil)

	if err := Interpret(module, namespace); err != nil {
		fmt.Println(exception.NewRuntimeError(fmt.Sprint(err)))
	}

	return *namespace, nil
}
