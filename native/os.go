package native

import (
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func readFile(args []interface{}) (interface{}, error) {
	path := args[0].(string)
	content, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewRuntimeError("RT032", "read", path)
	}

	return string(content[:]), nil
}

func writeFile(args []interface{}) (interface{}, error) {
	path := args[0].(string)
	data := args[1].(string)
	err := os.WriteFile(path, []byte(data), 0644)

	if err != nil {
		return nil, exception.NewRuntimeError("RT032", "write", path)
	}

	return nil, nil
}

func deleteFile(args []interface{}) (interface{}, error) {
	path := args[0].(string)
	err := os.Remove(path)

	if err != nil {
		return nil, exception.NewRuntimeError("RT032", "delete", path)
	}

	return nil, nil
}

var OsModule = InternalModule{
	symbols: map[string]InternalModuleFn{
		"readFile":   readFile,
		"writeFile":  writeFile,
		"deleteFile": deleteFile,
	},
}
