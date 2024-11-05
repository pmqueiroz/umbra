package helpers

import (
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func ReadFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewGenericError("GN001", path)
	}

	return string(dat[:]), nil
}

func UmbraType(value interface{}) string {
	switch value.(type) {
	case string:
		return "<str>"
	case bool:
		return "<bool>"
	case float64:
		return "<num>"
	case nil:
		return "<null>"
	case map[interface{}]interface{}:
		return "<hashmap>"
	case []interface{}:
		return "<arr>"
	default:
		return "<unknown>"
	}
}
