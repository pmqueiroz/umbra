package interpreter

import (
	"fmt"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/tokens"
)

func Eval(expression ast.Expression, env *Environment) (interface{}, error) {
	switch expr := expression.(type) {
	case ast.LiteralExpression:
		return expr.Value, nil

	case ast.VariableExpression:
		value, ok := env.Get(expr.Name.Raw.Value)
		if !ok {
			return nil, fmt.Errorf("variável não definida: %s", expr.Name.Raw.Value)
		}
		return value, nil

	case ast.AssignExpression:
		value, err := Eval(expr.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(expr.Name.Raw.Value, value)
		return value, nil

	case ast.BinaryExpression:
		left, err := Eval(expr.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := Eval(expr.Right, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Id {
		case tokens.PLUS:
			return left.(float64) + right.(float64), nil
		case tokens.MINUS:
			return left.(float64) - right.(float64), nil
		case tokens.STAR:
			return left.(float64) * right.(float64), nil
		case tokens.SLASH:
			if right.(float64) == 0 {
				return nil, fmt.Errorf("divisão por zero")
			}
			return left.(float64) / right.(float64), nil
		default:
			return nil, fmt.Errorf("operador desconhecido: %s", expr.Operator.Raw.Value)
		}

	case ast.UnaryExpression:
		right, err := Eval(expr.Right, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Id {
		case tokens.MINUS:
			return -right.(float64), nil
		default:
			return nil, fmt.Errorf("operador unário desconhecido: %s", expr.Operator.Raw.Value)
		}

	default:
		return nil, fmt.Errorf("expressão desconhecida: %T", expr)
	}
}

func Exec(stmt ast.Statement, env *Environment) error {
	switch s := stmt.(type) {
	case ast.PrintStatement:
		value, err := Eval(s.Expression, env)
		if err != nil {
			return err
		}
		fmt.Println(value)
		return nil
	case ast.VarStatement:
		var value interface{}
		var err error
		if s.Initializer != nil {
			value, err = Eval(s.Initializer, env)
			if err != nil {
				return err
			}
		}
		env.Set(s.Name.Raw.Value, value)
		return nil
	case ast.BlockStatement:
		newEnv := NewEnvironment(env)
		for _, stmt := range s.Statements {
			if err := Exec(stmt, newEnv); err != nil {
				return err
			}
		}
		return nil
	case ast.ModuleStatement:
		for _, stmt := range s.Declarations {
			if err := Exec(stmt, env); err != nil {
				return err
			}
		}
		return nil
	case ast.IfStatement:
		condition, err := Eval(s.Condition, env)
		if err != nil {
			return err
		}

		if condition.(bool) {
			return Exec(s.ThenBranch, env)
		} else if s.ElseBranch != nil {
			return Exec(s.ElseBranch, env)
		}
		return nil
	case ast.ReturnStatement:
		value, err := Eval(s.Value, env)
		if err != nil {
			return err
		}
		// Retornar o valor de uma função (simulação)
		return fmt.Errorf("return %v", value)
	default:
		return fmt.Errorf("declaração desconhecida: %T", stmt)
	}
}
