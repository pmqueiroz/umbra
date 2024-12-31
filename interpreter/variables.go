package interpreter

import (
	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
)

func resolveVarDeclaration(stmt ast.VarStatement, value interface{}, env *environment.Environment) error {
	err := checkDeclarationType(stmt.Type, stmt.Nullable, value, env, stmt)

	if err != nil {
		return err
	}

	varType, err := types.ParseTokenType(stmt.Type.Type)

	if err != nil {
		return err
	}

	env.Create(stmt, stmt.Name.Lexeme, value, varType, stmt.Nullable, false, stmt.Mutable)
	return nil
}

func zero(t tokens.TokenType) interface{} {
	switch t {
	case tokens.STR_TYPE:
		return ""
	case tokens.CHAR_TYPE:
		return rune(0)
	case tokens.BOOL_TYPE:
		return false
	case tokens.NUM_TYPE:
		return 0.0
	case tokens.HASHMAP_TYPE:
		return make(map[interface{}]interface{})
	case tokens.ARR_TYPE:
		return make([]interface{}, 0)
	case tokens.FUN_TYPE:
		return FunctionDeclaration{}
	default: // any, void
		return nil
	}
}
