package interpreter

import (
	"fmt"

	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
)

func CheckType(target tokens.TokenType, expected interface{}) error {
	typeMismatchErr := fmt.Sprintf("expected %s got %T", target, expected)

	switch expected.(type) {
	case nil:
		if target != tokens.NULL {
			return exception.NewTypeError(typeMismatchErr)
		}
	case string:
		if target != tokens.STR_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case bool:
		if target != tokens.BOOL_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case float64:
		if target != tokens.NUM_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case map[interface{}]interface{}:
		if target != tokens.HASHMAP_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case []interface{}:
		if target != tokens.ARR_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	}
	return nil
}
