package interpreter

import (
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/native"
	"github.com/pmqueiroz/umbra/tokens"
)

type Module struct {
	Name        string
	Environment *environment.Environment
}

func ResolveModule(module string) (string, error) {
	path := fmt.Sprintf("%s/lib/%s.u", os.Getenv("UMBRA_PATH"), module)

	content, err := helpers.ReadFile(path)

	if err != nil {
		return content, exception.NewGenericError("GN002", module, path)
	}

	return content, nil
}

func LoadInternalModule(name string, namespace *environment.Environment) error {
	if ok := native.Register(name, namespace); !ok {
		return fmt.Errorf("unable to include %s. internal module does not exits", name)
	}

	return nil
}

func LoadBuiltInModule(path string, namespace *environment.Environment) error {
	content, err := ResolveModule(path)

	if err != nil {
		return err
	}

	tokens, err := tokens.Tokenize(content)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	module := ast.Parse(tokens)

	if err := Interpret(module, namespace); err != nil {
		fmt.Println(err)
	}

	return nil
}

func LoadModule(path string) (Module, error) {
	namespace := environment.NewEnvironment(nil)

	if len(path) >= 7 && path[:7] == "native/" {
		err := LoadInternalModule(path[7:], namespace)

		if err != nil {
			return Module{}, err
		}

		return Module{
			Name:        path[7:],
			Environment: namespace,
		}, nil
	}

	err := LoadBuiltInModule(path, namespace)

	if err != nil {
		return Module{}, err
	}

	return Module{
		Name:        path,
		Environment: namespace,
	}, nil
}
