package types

import (
	"reflect"

	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/globals"
	"github.com/pmqueiroz/umbra/tokens"
)

func isFunctionDeclaration(value interface{}) bool {
	// Check if the value is an ast.FunctionDeclaration
	// TODO: Find a better way to do this preventing circular dependencies
	return reflect.TypeOf(value).Name() == "FunctionDeclaration"
}

func CheckPrimitiveType(targetType UmbraType, expected interface{}, nullable bool, node globals.Node) error {
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
	default:
		if isFunctionDeclaration(expected) {
			if targetType == FUN {
				return nil
			}
		}
	}

	expectedType, err := ParseUmbraType(expected)

	if err != nil {
		return exception.NewUmbraError("TY000", node, expected)
	}

	return exception.NewUmbraError("TY001", node, targetType, expectedType)
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
		if isFunctionDeclaration("FunctionDeclaration") {
			return FUN, nil
		}

		return UNKNOWN, exception.NewUmbraError("TY000", nil, value)
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
	case tokens.FUN_TYPE:
		return FUN, nil
	default:
		return UNKNOWN, exception.NewUmbraError("TY000", nil, value)
	}
}
