package interpreter

import (
	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/types"
)

func processFunction(funcExpr ast.FunctionExpression, env *environment.Environment) (FunctionDeclaration, error) {
	parsedReturnType, parentEnum, err := parseRuntimeType(funcExpr.ReturnType, env)

	if err != nil {
		return FunctionDeclaration{}, err
	}

	fun := FunctionDeclaration{Itself: &funcExpr, Environment: env, ReturnType: struct {
		Type   types.UmbraType
		Parent ast.EnumStatement
	}{Type: parsedReturnType, Parent: parentEnum}}

	if funcExpr.Name.Lexeme != "" {
		env.Create(
			funcExpr,
			funcExpr.Name.Lexeme,
			fun,
			types.FUN,
			false,
			false,
			false,
		)
	}

	return fun, nil
}

func resolveArgs(args interface{}, env *environment.Environment) ([]interface{}, error) {
	result := make([]interface{}, 0)

	switch args := args.(type) {
	case []ast.Expression:
		for _, arg := range args {
			value, err := Evaluate(arg, env)
			if err != nil {
				return nil, err
			}

			result = append(result, value)
		}
	case []ast.EnumArgument:
		for _, arg := range args {
			result = append(result, arg.Value)
		}
	}

	return result, nil
}

func processFunctionCall(callee FunctionDeclaration, args interface{}, env *environment.Environment) (interface{}, error) {
	funcEnv := environment.NewEnvironment(callee.Environment)
	parsedArgs, err := resolveArgs(args, env)

	if err != nil {
		return nil, err
	}

	for i, param := range callee.Itself.Params {
		if param.Variadic {
			var variadicArgs []interface{}
			for j := i; j < len(parsedArgs); j++ {
				typeErr := types.CheckPrimitiveType(param.Type, parsedArgs[j], param.Nullable, nil)
				if typeErr != nil {
					return nil, typeErr
				}

				variadicArgs = append(variadicArgs, parsedArgs[j])
			}
			funcEnv.Create(nil, param.Name.Lexeme, variadicArgs, param.Type, param.Nullable, false, false)
			break
		} else {
			typeErr := types.CheckPrimitiveType(param.Type, parsedArgs[i], param.Nullable, nil)
			if typeErr != nil {
				return nil, typeErr
			}

			funcEnv.Create(nil, param.Name.Lexeme, parsedArgs[i], param.Type, param.Nullable, false, false)
		}
	}

	var result interface{}
	for _, stmt := range callee.Itself.Body {
		if err := Interpret(stmt, funcEnv); err != nil {
			return nil, err
		}
	}

	return result, nil
}
