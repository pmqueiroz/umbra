package interpreter

import (
	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
)

func parseRuntimeType(t tokens.Token, env *environment.Environment) (types.UmbraType, ast.EnumStatement, error) {
	switch t.Type {
	case tokens.IDENTIFIER:
		value, ok := env.Get(t.Lexeme, true)
		if !ok {
			return types.UNKNOWN, ast.EnumStatement{}, exception.NewUmbraError("RT002", nil, t.Lexeme)
		}

		if enum, ok := value.Data.(ast.EnumStatement); ok {
			return types.ENUM, enum, nil
		}

		return types.UNKNOWN, ast.EnumStatement{}, exception.NewUmbraError("RT035", nil, t.Lexeme)
	default:
		t, e := types.ParseTokenType(t.Type)
		return t, ast.EnumStatement{}, e
	}
}
