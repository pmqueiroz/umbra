package interpreter

import "github.com/pmqueiroz/umbra/ast"

func checkMatch(pattern interface{}, expr interface{}) bool {
	switch p := pattern.(type) {
	case ast.EnumMember:
		if e, ok := expr.(ast.EnumMember); ok {
			// deep compare
			return p.Name == e.Name && p.Signature == e.Signature
		}

		return false
	default:
		return false
	}
}
