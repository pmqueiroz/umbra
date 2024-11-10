package types

import (
	"fmt"

	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
)

func CheckType(targetType tokens.TokenType, expected interface{}, nullable bool) error {
	typeMismatchErr := fmt.Sprintf("expected %s got %s", targetType, ParseUmbraType(expected))

	if targetType == tokens.ANY_TYPE {
		return nil
	}

	switch expected.(type) {
	case nil:
		if targetType != tokens.NULL && !nullable {
			return exception.NewTypeError(typeMismatchErr)
		}
	case string:
		if targetType != tokens.STR_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case rune:
		if targetType != tokens.CHAR_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case bool:
		if targetType != tokens.BOOL_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case float64:
		if targetType != tokens.NUM_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case map[interface{}]interface{}:
		if targetType != tokens.HASHMAP_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	case []interface{}:
		if targetType != tokens.ARR_TYPE {
			return exception.NewTypeError(typeMismatchErr)
		}
	}
	return nil
}

func ParseUmbraType(value interface{}) UmbraType {
	switch value.(type) {
	case string:
		return STR
	case rune:
		return CHAR
	case bool:
		return BOOL
	case float64:
		return NUM
	case nil:
		return NULL
	case map[interface{}]interface{}:
		return HASHMAP
	case []interface{}:
		return ARR
	default:
		return UNKNOWN
	}
}
