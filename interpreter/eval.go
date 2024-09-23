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
		value, ok := env.Get(expr.Name.Lexeme)
		if !ok {
			return nil, fmt.Errorf("undefined variable: %s", expr.Name.Lexeme)
		}
		return value, nil
	case ast.AssignExpression:
		value, err := Evaluate(expr.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(expr.Name.Lexeme, value)
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

		switch expr.Operator.Type {
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
			return nil, fmt.Errorf("unknown binary expression: %s", expr.Operator.Lexeme)
		}
	case ast.UnaryExpression:
		right, err := Evaluate(expr.Right, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Type {
		case tokens.MINUS:
			return -right.(float64), nil
		default:
			return nil, fmt.Errorf("unknown unary expression: %s", expr.Operator.Lexeme)
		}
	case ast.CallExpression:
		callee, err := Evaluate(expr.Callee, env)

		if err != nil {
			return nil, err
		}

		if function, ok := callee.(FunctionDeclaration); ok {
			funcEnv := NewEnvironment(function.Environment)

			for i, arg := range expr.Arguments {
				argValue, err := Evaluate(arg, env)
				if err != nil {
					return nil, err
				}
				funcEnv.Create(function.Itself.Params[i].Name.Lexeme, argValue)
			}

			var result interface{}
			for _, stmt := range function.Itself.Body {
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
	case ast.LogicalExpression:
		left, err := Evaluate(expr.Left, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Type {
		case tokens.OR:
			if left.(bool) {
				return true, nil
			}
		case tokens.AND:
			if !left.(bool) {
				return false, nil
			}
		default:
			return nil, fmt.Errorf("unknown logical operator: %s", expr.Operator.Lexeme)
		}

		right, err := Evaluate(expr.Right, env)
		if err != nil {
			return nil, err
		}

		return right, nil
	case ast.ArrayExpression:
		var elements []interface{}
		for _, element := range expr.Elements {
			evaluatedElement, err := Evaluate(element, env)
			if err != nil {
				return nil, err
			}
			elements = append(elements, evaluatedElement)
		}
		return elements, nil
	case ast.HashmapExpression:
		hashmap := make(map[interface{}]interface{})
		for keyExpr, valueExpr := range expr.Pairs {
			key, err := Evaluate(keyExpr, env)
			if err != nil {
				return nil, err
			}
			value, err := Evaluate(valueExpr, env)
			if err != nil {
				return nil, err
			}
			hashmap[key] = value
		}
		return hashmap, nil
	case ast.MemberExpression:
		object, err := Evaluate(expr.Object, env)
		if err != nil {
			return nil, err
		}

		switch obj := object.(type) {
		case map[interface{}]interface{}:
			value, ok := obj[expr.Property.Lexeme]
			if !ok {
				return nil, fmt.Errorf("undefined property: %v", expr.Property)
			}
			return value, nil
		case []interface{}:
			index, err := Evaluate(expr.Property, env)
			if err != nil {
				return nil, err
			}
			idx, ok := index.(float64)
			if !ok {
				return nil, fmt.Errorf("invalid array index: %v", index)
			}
			if int(idx) < 0 || int(idx) >= len(obj) {
				return nil, fmt.Errorf("array index out of bounds: %v", idx)
			}
			return obj[int(idx)], nil
		default:
			return nil, fmt.Errorf("cannot access property of non-object type: %T", obj)
		}
	default:
		return nil, fmt.Errorf("unknown expression: %T", expr)
	}
}
