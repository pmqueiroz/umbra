package interpreter

import (
	"fmt"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/tokens"
)

func Evaluate(expression ast.Expression, env *Environment) (interface{}, error) {
	switch expr := expression.(type) {
	case ast.LiteralExpression:
		return expr.Value, nil
	case ast.GroupingExpression:
		return Evaluate(expr.Expression, env)
	case ast.VariableExpression:
		value, ok := env.Get(expr.Name.Raw.Value)
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", expr.Name.Raw.Value)
		}
		return value, nil
	case ast.AssignExpression:
		value, err := Evaluate(expr.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(expr.Name.Raw.Value, value)
		return value, nil
	case ast.BinaryExpression:
		left, err := Evaluate(expr.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := Evaluate(expr.Right, env)
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
				return nil, fmt.Errorf("invalid operation: division by zero")
			}
			return left.(float64) / right.(float64), nil
		case tokens.GREATER_THAN:
			return left.(float64) > right.(float64), nil
		case tokens.GREATER_THAN_EQUAL:
			return left.(float64) >= right.(float64), nil
		case tokens.LESS_THAN:
			return left.(float64) < right.(float64), nil
		case tokens.LESS_THAN_EQUAL:
			return left.(float64) <= right.(float64), nil
		case tokens.EQUAL_EQUAL:
			return left == right, nil
		case tokens.BANG_EQUAL:
			return left != right, nil
		default:
			return nil, fmt.Errorf("unknown binary expression: %s", expr.Operator.Raw.Value)
		}
	case ast.UnaryExpression:
		right, err := Evaluate(expr.Right, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Id {
		case tokens.MINUS:
			return -right.(float64), nil
		default:
			return nil, fmt.Errorf("unknown unary expression: %s", expr.Operator.Raw.Value)
		}
	case ast.CallExpression:
		callee, err := Evaluate(expr.Callee, env)

		if err != nil {
			return nil, err
		}

		if function, ok := callee.(ast.FunctionStatement); ok {
			funcEnv := NewEnvironment(env)

			for i, arg := range expr.Arguments {
				argValue, err := Evaluate(arg, env)
				if err != nil {
					return nil, err
				}
				funcEnv.Set(function.Params[i].Name.Raw.Value, argValue)
			}

			var result interface{}
			for _, stmt := range function.Body {
				if err := Interpret(stmt, funcEnv); err != nil {
					if returnValue, ok := err.(Return); ok {
						result = returnValue.value
						break
					}

					return nil, err
				}
			}

			return result, nil
		}

		return nil, fmt.Errorf("invalid function call %v", expr.Callee)
	default:
		return nil, fmt.Errorf("unknown expression: %T", expr)
	}
}
