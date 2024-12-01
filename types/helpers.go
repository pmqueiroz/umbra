package types

import (
	"fmt"

	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
)

func CheckPrimitiveType(targetType UmbraType, expected interface{}, nullable bool) error {
	if targetType == ANY {
		return nil
	}

	switch expected.(type) {
	case nil:
		if targetType == NULL || nullable {
			return nil
		}
	case string:
		if targetType == STR {
			return nil
		}
	case rune:
		if targetType == CHAR {
			return nil
		}
	case bool:
		if targetType == BOOL {
			return nil
		}
	case float64:
		if targetType == NUM {
			return nil
		}
	case map[interface{}]interface{}:
		if targetType == HASHMAP {
			return nil
		}
	case []interface{}:
		if targetType == ARR {
			return nil
		}
	}

	expectedType, err := ParseUmbraType(expected)

	if err != nil {
		return exception.NewTypeError(fmt.Sprintf("type %s is invalid", expected))
	}

	return exception.NewTypeError(fmt.Sprintf("expected %s got %s", targetType, expectedType))
}

func ParseUmbraType(value interface{}) (UmbraType, error) {
	switch value.(type) {
	case string:
		return STR, nil
	case rune:
		return CHAR, nil
	case bool:
		return BOOL, nil
	case float64:
		return NUM, nil
	case nil:
		return NULL, nil
	case map[interface{}]interface{}:
		return HASHMAP, nil
	case []interface{}:
		return ARR, nil
	default:
		return UNKNOWN, exception.NewTypeError("unknown type")
	}
}

func SafeParseUmbraType(value interface{}) UmbraType {
	umbraType, err := ParseUmbraType(value)
	if err != nil {
		return UNKNOWN
	}
	return umbraType
}

func ParseTokenType(value tokens.TokenType) (UmbraType, error) {
	switch value {
	case tokens.STR_TYPE:
		return STR, nil
	case tokens.CHAR_TYPE:
		return CHAR, nil
	case tokens.BOOL_TYPE:
		return BOOL, nil
	case tokens.NUM_TYPE:
		return NUM, nil
	case tokens.HASHMAP_TYPE:
		return HASHMAP, nil
	case tokens.ARR_TYPE:
		return ARR, nil
	case tokens.ANY_TYPE:
		return ANY, nil
	case tokens.VOID_TYPE:
		return VOID, nil
	default:
		return UNKNOWN, exception.NewTypeError("unknown type")
	}
}
