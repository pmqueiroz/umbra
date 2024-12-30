package native

import (
	"path/filepath"

	"github.com/pmqueiroz/umbra/exception"
)

func allString(maybeStrings []interface{}) ([]string, error) {
	strPaths := make([]string, len(maybeStrings))
	for i, p := range maybeStrings {
		str, ok := p.(string)
		if !ok {
			return nil, exception.NewRuntimeError("RT033", "")
		}
		strPaths[i] = str
	}

	return strPaths, nil
}

func resolve(paths []interface{}) (interface{}, error) {
	strPaths, err := allString(paths)
	if err != nil {
		return nil, err
	}

	return filepath.Join(strPaths...), nil
}

func dirname(args []interface{}) (interface{}, error) {
	path, ok := args[0].(string)
	if !ok {
		return nil, exception.NewRuntimeError("RT033", "")
	}

	return filepath.Dir(path), nil
}

var PathModule = InternalModule{
	symbols: map[string]InternalModuleFn{
		"resolve": resolve,
		"dirname": dirname,
	},
}
