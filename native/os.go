package native

import (
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func readFile(args []interface{}) (interface{}, error) {
	path := args[0].(string)
	content, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewRuntimeError("RT032", path)
	}

	return string(content[:]), nil
}

var OsModule = InternalModule{
	symbols: map[string]InternalModuleFn{
		"ReadFile": readFile,
	},
}
