package interpreter

import (
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/native"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
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

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case rune:
		return float64(v), nil
	default:
		return 0, exception.NewRuntimeError("RT026", types.SafeParseUmbraType(value))
	}
}

func Evaluate(expression ast.Expression, env *environment.Environment) (interface{}, error) {
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
		return value.Data, nil
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

			typeErr := types.CheckPrimitiveType(variable.DataType, value, variable.Nullable)

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
				return nil, exception.NewRuntimeError("RT005", types.SafeParseUmbraType(obj))
			}
		default:
			return nil, exception.NewRuntimeError("RT006", types.SafeParseUmbraType(target))
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
			switch leftVal := left.(type) {
			case string:
				switch rightVal := right.(type) {
				case string:
					return leftVal + rightVal, nil
				case rune:
					return leftVal + string(rightVal), nil
				}
			case rune:
				if rightStr, ok := right.(string); ok {
					return string(leftVal) + rightStr, nil
				} else if rightRune, ok := right.(rune); ok {
					return leftVal + rightRune, nil
				} else if rightFloat, ok := right.(float64); ok {
					return leftVal + rune(rightFloat), nil
				}
			case float64:
				if rightFloat, ok := right.(float64); ok {
					return leftVal + rightFloat, nil
				}
			}
			return nil, exception.NewRuntimeError("RT007", types.SafeParseUmbraType(left), types.SafeParseUmbraType(right))
		case tokens.MINUS:
			switch leftVal := left.(type) {
			case float64:
				if rightFloat, ok := right.(float64); ok {
					return leftVal - rightFloat, nil
				} else if rightRune, ok := right.(rune); ok {
					return leftVal - float64(rightRune), nil
				}
			case rune:
				if rightFloat, ok := right.(float64); ok {
					return rune(float64(leftVal) - rightFloat), nil
				} else if rightRune, ok := right.(rune); ok {
					return leftVal - rightRune, nil
				}
			}
			return nil, exception.NewRuntimeError("RT027", types.SafeParseUmbraType(left), types.SafeParseUmbraType(right))
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

			return nil, exception.NewRuntimeError("RT009", types.SafeParseUmbraType(left), types.SafeParseUmbraType(right))
		case tokens.GREATER_THAN:
			leftVal, err := toFloat64(left)
			if err != nil {
				return false, err
			}
			rightVal, err := toFloat64(right)
			if err != nil {
				return false, err
			}
			return leftVal > rightVal, nil

		case tokens.GREATER_THAN_EQUAL:
			leftVal, err := toFloat64(left)
			if err != nil {
				return false, err
			}
			rightVal, err := toFloat64(right)
			if err != nil {
				return false, err
			}
			return leftVal >= rightVal, nil

		case tokens.LESS_THAN:
			leftVal, err := toFloat64(left)
			if err != nil {
				return false, err
			}
			rightVal, err := toFloat64(right)
			if err != nil {
				return false, err
			}
			return leftVal < rightVal, nil

		case tokens.LESS_THAN_EQUAL:
			leftVal, err := toFloat64(left)
			if err != nil {
				return false, err
			}
			rightVal, err := toFloat64(right)
			if err != nil {
				return false, err
			}
			return leftVal <= rightVal, nil
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
		case tokens.TYPE_OF:
			parsedType, err := types.ParseUmbraType(right)

			if err != nil {
				return nil, err
			}

			return parsedType, nil
		case tokens.TILDE:
			switch parsedRight := right.(type) {
			case []interface{}:
				return float64(len(parsedRight)), nil
			case []string:
				return float64(len(parsedRight)), nil
			case string:
				return float64(len(parsedRight)), nil
			default:
				return nil, exception.NewRuntimeError("RT011", types.SafeParseUmbraType(parsedRight))
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
				return nil, exception.NewRuntimeError("RT012", types.SafeParseUmbraType(parsedRight))
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
			funcEnv := environment.NewEnvironment(function.Environment)

			for i, param := range function.Itself.Params {
				if param.Variadic {
					var variadicArgs []interface{}
					for j := i; j < len(expr.Arguments); j++ {
						argValue, err := Evaluate(expr.Arguments[j], env)
						if err != nil {
							return nil, err
						}

						typeErr := types.CheckPrimitiveType(param.Type, argValue, param.Nullable)
						if typeErr != nil {
							return nil, typeErr
						}

						variadicArgs = append(variadicArgs, argValue)
					}
					funcEnv.Create(param.Name.Lexeme, variadicArgs, param.Type, param.Nullable, false)
					break
				} else {
					argValue, err := Evaluate(expr.Arguments[i], env)
					if err != nil {
						return nil, err
					}

					typeErr := types.CheckPrimitiveType(param.Type, argValue, param.Nullable)
					if typeErr != nil {
						return nil, typeErr
					}

					funcEnv.Create(param.Name.Lexeme, argValue, param.Type, param.Nullable, false)
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
		} else if goFunc, ok := callee.(native.InternalModuleFn); ok {
			var args []interface{}
			for _, arg := range expr.Arguments {
				argValue, err := Evaluate(arg, env)
				if err != nil {
					return nil, err
				}
				args = append(args, argValue)
			}

			defer func() {
				if r := recover(); r != nil {
					err = exception.NewRuntimeError("RT031")
				}
			}()
			result, err := goFunc(args)
			return result, err
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
		case string:
			index, err := Evaluate(expr.Property, env)
			if err != nil {
				return nil, err
			}
			idx, ok := index.(float64)
			if !ok {
				return nil, exception.NewRuntimeError("RT003", index)
			}

			value := []rune(obj)[int(idx)]

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
		case ast.EnumStatement:
			if prop, ok := expr.Property.(ast.VariableExpression); ok {
				member, ok := obj.Get(prop.Name)
				if !ok {
					return nil, exception.NewRuntimeError("RT034", prop.Name.Lexeme)
				}

				return member, nil
			}
			return nil, exception.NewRuntimeError("RT020")
		default:
			return nil, exception.NewRuntimeError("RT016", types.SafeParseUmbraType(obj))
		}
	case ast.NamespaceMemberExpression:
		if variableExpr, ok := expr.Namespace.(ast.VariableExpression); ok {
			namespace, ok := env.GetNamespace(variableExpr.Name.Lexeme)
			if !ok {
				return nil, exception.NewRuntimeError("RT018", variableExpr.Name.Lexeme)
			}

			value, _ := namespace.Get(expr.Property.Lexeme, false)

			return value.Data, nil
		}

		return nil, exception.NewRuntimeError("RT019", litter.Sdump(expr.Namespace))
	case ast.TypeConversionExpression:
		value, err := Evaluate(expr.Value, env)
		if err != nil {
			return nil, err
		}

		defaultError := exception.NewRuntimeError("RT028", types.SafeParseUmbraType(value), types.SafeParseUmbraType(expr.Type.Type))

		switch expr.Type.Type {
		case tokens.STR_TYPE:
			switch v := value.(type) {
			case rune:
				return string(v), nil
			case float64:
				return strconv.FormatFloat(v, 'f', -1, 64), nil
			case bool:
				return strconv.FormatBool(v), nil
			}
			return nil, defaultError
		case tokens.CHAR_TYPE:
			switch v := value.(type) {
			case float64:
				return rune(v), nil
			case string:
				switch utf8.RuneCountInString(v) {
				case 1:
					return []rune(v)[0], nil
				case 2:
					if []rune(v)[0] != '\\' {
						return nil, exception.NewRuntimeError("RT029")
					}

					runeStr, err := strconv.Unquote(`"` + v + `"`)
					if err != nil {
						return nil, exception.NewRuntimeError("RT030")
					}

					return []rune(runeStr)[0], nil
				default:
					return nil, exception.NewRuntimeError("RT029")
				}
			}
			return nil, defaultError
		case tokens.NUM_TYPE:
			switch v := value.(type) {
			case bool:
				return map[bool]float64{true: 1.0, false: 0.0}[v], nil
			case rune:
				return float64(v), nil
			case string:
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return math.NaN(), nil
				}

				return value, nil
			}
			return math.NaN(), defaultError
		}
		return nil, defaultError
	default:
		return nil, exception.NewRuntimeError("RT017", litter.Sdump(expr))
	}
}

func resolveMemberExpressionProperty(expr ast.MemberExpression, env *environment.Environment) (interface{}, error) {
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
