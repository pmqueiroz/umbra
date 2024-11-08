package interpreter

import (
	"math"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/sanity-io/litter"
)

func getLength(data interface{}) int {
	switch v := data.(type) {
	case []string:
		return len(v)
	case []interface{}:
		return len(v)
	default:
		return -1
	}
}

func getElementAt(data interface{}, idx int) interface{} {
	switch v := data.(type) {
	case []string:
		return v[idx]
	case []interface{}:
		return v[idx]
	default:
		return nil
	}
}

func Evaluate(expression ast.Expression, env *Environment) (interface{}, error) {
	switch expr := expression.(type) {
	case ast.LiteralExpression:
		return expr.Value, nil
	case ast.NaNExpression:
		return math.NaN(), nil
	case ast.GroupingExpression:
		return Evaluate(expr.Expression, env)
	case ast.VariableExpression:
		value, ok := env.Get(expr.Name.Lexeme, true)
		if !ok {
			return nil, exception.NewRuntimeError("RT002", expr.Name.Lexeme)
		}
		return value.data, nil
	case ast.AssignExpression:
		value, err := Evaluate(expr.Value, env)
		if err != nil {
			return nil, err
		}

		switch target := expr.Target.(type) {
		case ast.VariableExpression:
			variable, exists := env.Get(target.Name.Lexeme, true)

			if !exists {
				return nil, exception.NewRuntimeError("RT002", target.Name.Lexeme)
			}

			typeErr := CheckType(variable.dataType, value, variable.nullable)

			if typeErr != nil {
				return nil, typeErr
			}

			env.Set(target.Name.Lexeme, value)
			return value, nil
		case ast.MemberExpression:
			object, err := Evaluate(target.Object, env)
			if err != nil {
				return nil, err
			}

			property, err := resolveMemberExpressionProperty(target, env)

			if err != nil {
				return nil, err
			}

			switch obj := object.(type) {
			case map[interface{}]interface{}:
				obj[property] = value
				return value, nil
			case []interface{}:
				index, err := Evaluate(target.Property, env)
				if err != nil {
					return nil, err
				}
				idx, ok := index.(float64)
				if !ok {
					return nil, exception.NewRuntimeError("RT003", index)
				}
				if int(idx) < 0 || int(idx) > len(obj) {
					return nil, exception.NewRuntimeError("RT004", idx)
				}
				if int(idx) == len(obj) {
					env.Set(target.Object.(ast.VariableExpression).Name.Lexeme, append(obj, value))
					return value, nil
				}
				obj[int(idx)] = value
				return value, nil
			default:
				return nil, exception.NewRuntimeError("RT005", helpers.UmbraType(obj))
			}
		default:
			return nil, exception.NewRuntimeError("RT006", helpers.UmbraType(target))
		}
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
			leftStr, leftIsString := left.(string)
			rightStr, rightIsString := right.(string)
			if leftIsString && rightIsString {
				return leftStr + rightStr, nil
			}
			leftFloat, leftIsFloat := left.(float64)
			rightFloat, rightIsFloat := right.(float64)
			if leftIsFloat && rightIsFloat {
				return leftFloat + rightFloat, nil
			}
			return nil, exception.NewRuntimeError("RT007", helpers.UmbraType(left), helpers.UmbraType(right))
		case tokens.MINUS:
			return left.(float64) - right.(float64), nil
		case tokens.STAR:
			return left.(float64) * right.(float64), nil
		case tokens.SLASH:
			if right.(float64) == 0 {
				return nil, exception.NewRuntimeError("RT008")
			}
			return left.(float64) / right.(float64), nil
		case tokens.PERCENT:
			leftFloat, leftIsFloat := left.(float64)
			rightFloat, rightIsFloat := right.(float64)
			if leftIsFloat && rightIsFloat {
				return math.Mod(leftFloat, rightFloat), nil
			}

			return nil, exception.NewRuntimeError("RT009", helpers.UmbraType(left), helpers.UmbraType(right))
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
			return nil, exception.NewRuntimeError("RT010", expr.Operator.Lexeme)
		}
	case ast.UnaryExpression:
		right, err := Evaluate(expr.Right, env)
		if err != nil {
			return nil, err
		}

		switch expr.Operator.Type {
		case tokens.MINUS:
			return -right.(float64), nil
		case tokens.NOT:
			return !right.(bool), nil
		case tokens.TILDE:
			switch parsedRight := right.(type) {
			case []interface{}:
				return float64(len(parsedRight)), nil
			case []string:
				return float64(len(parsedRight)), nil
			case string:
				return float64(len(parsedRight)), nil
			default:
				return nil, exception.NewRuntimeError("RT011", helpers.UmbraType(parsedRight))
			}
		case tokens.RANGE:
			switch parsedRight := right.(type) {
			case string:
				return strings.Split(parsedRight, ""), nil
			case map[interface{}]interface{}:
				var result [][]interface{}
				for key, value := range parsedRight {
					result = append(result, []interface{}{key, value})
				}
				return result, nil
			default:
				return nil, exception.NewRuntimeError("RT012", helpers.UmbraType(parsedRight))
			}
		default:
			return nil, exception.NewRuntimeError("RT013", expr.Operator.Lexeme)
		}
	case ast.CallExpression:
		callee, err := Evaluate(expr.Callee, env)

		if err != nil {
			return nil, err
		}

		if function, ok := callee.(FunctionDeclaration); ok {
			funcEnv := NewEnvironment(function.Environment)

			for i, param := range function.Itself.Params {
				if param.Variadic {
					var variadicArgs []interface{}
					for j := i; j < len(expr.Arguments); j++ {
						argValue, err := Evaluate(expr.Arguments[j], env)
						if err != nil {
							return nil, err
						}

						typeErr := CheckType(param.Type.Type, argValue, param.Nullable)
						if typeErr != nil {
							return nil, typeErr
						}

						variadicArgs = append(variadicArgs, argValue)
					}
					funcEnv.Create(param.Name.Lexeme, variadicArgs, param.Type.Type, param.Nullable)
					break
				} else {
					argValue, err := Evaluate(expr.Arguments[i], env)
					if err != nil {
						return nil, err
					}

					typeErr := CheckType(param.Type.Type, argValue, param.Nullable)
					if typeErr != nil {
						return nil, typeErr
					}

					funcEnv.Create(param.Name.Lexeme, argValue, param.Type.Type, param.Nullable)
				}
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

		return nil, exception.NewRuntimeError("RT014", expr.Callee)
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
			return nil, exception.NewRuntimeError("RT015", expr.Operator.Lexeme)
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

		property, err := resolveMemberExpressionProperty(expr, env)

		if err != nil {
			return nil, err
		}

		switch obj := object.(type) {
		case map[interface{}]interface{}:
			value, ok := obj[property]
			if !ok {
				return nil, nil
			}
			return value, nil
		case []interface{}, []string:
			index, err := Evaluate(expr.Property, env)
			if err != nil {
				return nil, err
			}
			idx, ok := index.(float64)
			if !ok {
				return nil, exception.NewRuntimeError("RT003", index)
			}
			if int(idx) < 0 || int(idx) >= getLength(obj) {
				return nil, exception.NewRuntimeError("RT004", idx)
			}
			return getElementAt(obj, int(idx)), nil
		default:
			return nil, exception.NewRuntimeError("RT016", helpers.UmbraType(obj))
		}
	case ast.NamespaceMemberExpression:
		if variableExpr, ok := expr.Namespace.(ast.VariableExpression); ok {
			namespace, ok := env.GetNamespace(variableExpr.Name.Lexeme)
			if !ok {
				return nil, exception.NewRuntimeError("RT018", variableExpr.Name.Lexeme)
			}

			value, _ := namespace.Get(expr.Property.Lexeme, false)

			return value.data, nil
		}

		return nil, exception.NewRuntimeError("RT019", litter.Sdump(expr.Namespace))
	default:
		return nil, exception.NewRuntimeError("RT017", litter.Sdump(expr))
	}
}

func resolveMemberExpressionProperty(expr ast.MemberExpression, env *Environment) (interface{}, error) {
	var property interface{}
	var computeErr error
	if expr.Computed {
		property, computeErr = Evaluate(expr.Property, env)
	} else if variable, ok := expr.Property.(ast.VariableExpression); ok {
		property = variable.Name.Lexeme
	} else {
		return nil, exception.NewRuntimeError("RT020")
	}

	if computeErr != nil {
		return nil, computeErr
	}

	return property, nil
}
