package interpreter

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/environment"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
	"github.com/sanity-io/litter"
)

type Return struct {
	value interface{}
}

func (r Return) Error() string {
	return "function returned"
}

type Break struct{}

func (r Break) Error() string {
	return "for loop break"
}

type Continue struct{}

func (r Continue) Error() string {
	return "for loop continue"
}

type FunctionDeclaration struct {
	Itself      *ast.FunctionStatement
	Environment *environment.Environment
	// temp solution while UmbraType is only a string
	ReturnType struct {
		Type   types.UmbraType
		Parent ast.EnumStatement
	}
}

func extractVarName(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case ast.VarStatement:
		return s.Name.Lexeme
	default:
		return ""
	}
}

func checkDeclarationType(t tokens.Token, nullable bool, value interface{}, env *environment.Environment) error {
	parsedType, enum, err := parseRuntimeType(t, env)

	if err != nil {
		return err
	}

	switch parsedType {
	case types.ENUM:
		if member, ok := value.(ast.EnumMember); ok {
			if ok := enum.ValidMember(member); !ok {
				return exception.NewTypeError(fmt.Sprintf("expected %s got %s", enum.Name.Lexeme, value))
			}
		}
	default:
		typeErr := types.CheckPrimitiveType(parsedType, value, nullable)

		if typeErr != nil {
			return typeErr
		}
	}

	return nil
}

func resolveVarDeclaration(stmt ast.VarStatement, value interface{}, env *environment.Environment) error {
	err := checkDeclarationType(stmt.Type, stmt.Nullable, value, env)

	if err != nil {
		return err
	}

	varType, err := types.ParseTokenType(stmt.Type.Type)

	if err != nil {
		return err
	}

	env.Create(stmt.Name.Lexeme, value, varType, stmt.Nullable, false)
	return nil
}

func Interpret(statement ast.Statement, env *environment.Environment) error {
	switch stmt := statement.(type) {
	case ast.PrintStatement:
		value, err := Evaluate(stmt.Expression, env)
		if err != nil {
			return err
		}
		var channel *os.File
		if stmt.Channel == ast.StderrChannel {
			channel = os.Stderr
		} else {
			channel = os.Stdout
		}

		var output string

		switch v := value.(type) {
		case string:
			v = strings.ReplaceAll(v, "\\n", "\n")
			v = strings.ReplaceAll(v, "\\t", "\t")
			v = strings.ReplaceAll(v, "\\\"", "\"")
			v = strings.ReplaceAll(v, "\\\\", "\\")
			output = fmt.Sprint(v)
		case float64:
			output = fmt.Sprintf("%.f", v)
		default:
			output = fmt.Sprint(value)
		}
		channel.Write([]byte(output))
		return nil
	case ast.VarStatement:
		var value interface{}
		var err error
		if stmt.Initializer != nil {
			value, err = Evaluate(stmt.Initializer, env)
			if err != nil {
				return err
			}
		}

		return resolveVarDeclaration(stmt, value, env)
	case ast.ArrayDestructuringStatement:
		value, err := Evaluate(stmt.Expr, env)
		if err != nil {
			return err
		}

		if result, ok := value.([]interface{}); ok {
			for i, declaration := range stmt.Declarations {
				err := resolveVarDeclaration(declaration, result[i], env)

				if err != nil {
					return err
				}
			}

			return nil
		}

		return exception.NewRuntimeError("RT039", types.SafeParseUmbraType(value))
	case ast.BlockStatement:
		newEnv := environment.NewEnvironment(env)
		for _, stmt := range stmt.Statements {
			if err := Interpret(stmt, newEnv); err != nil {
				return err
			}
		}
		return nil
	case ast.ModuleStatement:
		for _, stmt := range stmt.Declarations {
			if err := Interpret(stmt, env); err != nil {
				return err
			}
		}
		return nil
	case ast.IfStatement:
		condition, err := Evaluate(stmt.Condition, env)
		if err != nil {
			return err
		}

		if condition.(bool) {
			return Interpret(stmt.ThenBranch, env)
		} else if stmt.ElseBranch != nil {
			return Interpret(stmt.ElseBranch, env)
		}
		return nil
	case ast.ReturnStatement:
		value, err := Evaluate(stmt.Value, env)
		if err != nil {
			return err
		}
		return Return{value: value}
	case ast.FunctionStatement:
		parsedReturnType, parentEnum, err := parseRuntimeType(stmt.ReturnType, env)

		if err != nil {
			return err
		}

		env.Create(
			stmt.Name.Lexeme,
			FunctionDeclaration{Itself: &stmt, Environment: env, ReturnType: struct {
				Type   types.UmbraType
				Parent ast.EnumStatement
			}{Type: parsedReturnType, Parent: parentEnum}},
			types.FUN,
			false,
			false,
		)
		return nil
	case ast.ExpressionStatement:
		_, err := Evaluate(stmt.Expression, env)
		return err
	case ast.InitializedForStatement:
		forEnv := environment.NewEnvironment(env)
		if err := Interpret(stmt.Start, forEnv); err != nil {
			return err
		}

		initializedVarName := extractVarName(stmt.Start)

		for {
			loopEnv := environment.NewEnvironment(forEnv)
			controlVar, exists := loopEnv.Get(initializedVarName, true)
			if !exists {
				return exception.NewRuntimeError("RT021", initializedVarName)
			}

			stop, err := Evaluate(stmt.Stop, loopEnv)
			if err != nil {
				return err
			}

			var condition bool
			if parsedStop, ok := stop.(float64); ok {
				condition = controlVar.Data.(float64) <= parsedStop
			} else {
				return exception.NewRuntimeError("RT022", types.SafeParseUmbraType(stop))
			}

			if !condition {
				break
			}

			bodyErr := Interpret(stmt.Body, loopEnv)

			stepValue, err := Evaluate(stmt.Step, loopEnv)
			if err != nil {
				return err
			}

			step, exists := stepValue.(float64)
			if !exists {
				return exception.NewRuntimeError("RT023", types.SafeParseUmbraType(stepValue))
			}

			loopEnv.Set(initializedVarName, controlVar.Data.(float64)+step)

			if _, ok := bodyErr.(Break); ok {
				break
			}

			if _, ok := bodyErr.(Continue); ok {
				continue
			}

			if bodyErr != nil {
				return bodyErr
			}
		}

		return nil
	case ast.ConditionalForStatement:
		for {
			loopEnv := environment.NewEnvironment(env)

			condition, err := Evaluate(stmt.Condition, loopEnv)
			if err != nil {
				return err
			}

			parsedCondition, ok := condition.(bool)
			if !ok {
				return exception.NewRuntimeError("RT024", types.SafeParseUmbraType(parsedCondition))
			}

			if !parsedCondition {
				break
			}

			if err := Interpret(stmt.Body, loopEnv); err != nil {
				if _, ok := err.(Break); ok {
					break
				}

				if _, ok := err.(Continue); ok {
					continue
				}
				return err
			}
		}
		return nil
	case ast.BreakStatement:
		return Break{}
	case ast.ContinueStatement:
		return Continue{}
	case ast.PublicStatement:
		for _, identifier := range stmt.Identifiers {
			success := env.MakePublic(identifier.Lexeme)

			if !success {
				return exception.NewRuntimeError("RT025", identifier.Lexeme)
			}

		}
		return nil
	case ast.ImportStatement:
		module, err := LoadModule(stmt.Path.Lexeme)
		if err != nil {
			return err
		}
		env.CreateNamespace(module.Name, module.Environment)
		return nil
	case ast.EnumStatement:
		hasher := sha256.New()
		hasher.Write([]byte(stmt.Name.Lexeme))
		for member := range stmt.Members {
			hasher.Write([]byte(member))
		}

		stmt.Signature = hex.EncodeToString(hasher.Sum(nil))

		for name, member := range stmt.Members {
			member.Signature = stmt.Signature
			stmt.Members[name] = member
		}

		env.Create(
			stmt.Name.Lexeme,
			stmt,
			types.ENUM,
			false,
			false,
		)
		return nil
	case ast.MatchStatement:
		value, err := Evaluate(stmt.Expression, env)
		if err != nil {
			return err
		}

		for _, matchCase := range stmt.Cases {
			expr, err := Evaluate(matchCase.Expression, env)

			if err != nil {
				return err
			}

			if checkMatch(expr, value) {
				caseEnv := environment.NewEnvironment(env)

				for i, param := range matchCase.Parameters {
					arg := value.(ast.EnumMember).Arguments[i]

					caseEnv.Create(param.Name.Lexeme, arg.Value, types.ANY, false, false)
				}

				for _, stmt := range matchCase.Body {
					err := Interpret(stmt, caseEnv)

					if err != nil {
						return err
					}
				}

				return nil
			}
		}

		return nil
	default:
		litter.Dump("statement", statement)
		return exception.NewRuntimeError("RT000", reflect.TypeOf(statement).Name())
	}
}
