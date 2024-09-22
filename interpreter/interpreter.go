package interpreter

import (
	"fmt"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
)

type Return struct {
	value interface{}
}

func (r Return) Error() string {
	return "function returned"
}

func Interpret(stmt ast.Statement, env *Environment) error {
	switch s := stmt.(type) {
	case ast.PrintStatement:
		value, err := Evaluate(s.Expression, env)
		if err != nil {
			return err
		}
		if str, ok := value.(string); ok {
			str = strings.ReplaceAll(str, "\\n", "\n")
			str = strings.ReplaceAll(str, "\\t", "\t")
			str = strings.ReplaceAll(str, "\\\"", "\"")
			str = strings.ReplaceAll(str, "\\\\", "\\")
			fmt.Print(str)
		} else {
			fmt.Print(value)
		}
		return nil
	case ast.VarStatement:
		var value interface{}
		var err error
		if s.Initializer != nil {
			value, err = Evaluate(s.Initializer, env)
			if err != nil {
				return err
			}
		}
		env.Set(s.Name.Raw.Value, value)
		return nil
	case ast.BlockStatement:
		newEnv := NewEnvironment(env)
		for _, stmt := range s.Statements {
			if err := Interpret(stmt, newEnv); err != nil {
				return err
			}
		}
		return nil
	case ast.ModuleStatement:
		for _, stmt := range s.Declarations {
			if err := Interpret(stmt, env); err != nil {
				return err
			}
		}
		return nil
	case ast.IfStatement:
		condition, err := Evaluate(s.Condition, env)
		if err != nil {
			return err
		}

		if condition.(bool) {
			return Interpret(s.ThenBranch, env)
		} else if s.ElseBranch != nil {
			return Interpret(s.ElseBranch, env)
		}
		return nil
	case ast.ReturnStatement:
		value, err := Evaluate(s.Value, env)
		if err != nil {
			return err
		}
		return Return{value: value}
	case ast.FunctionStatement:
		env.Set(s.Name.Raw.Value, s)
		return nil
	case ast.ExpressionStatement:
		_, err := Evaluate(s.Expression, env)
		return err
	default:
		return fmt.Errorf("unknown declaration: %T", stmt)
	}
}
