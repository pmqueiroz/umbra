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
	"github.com/pmqueiroz/umbra/globals"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
)

type Return struct {
	Value interface{}
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
	Itself      *ast.FunctionExpression
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

func checkDeclarationType(t tokens.Token, nullable bool, value interface{}, env *environment.Environment, node globals.Node) error {
	parsedType, enum, err := parseRuntimeType(t, env)

	if err != nil {
		return err
	}

	switch parsedType {
	case types.ENUM:
		if member, ok := value.(ast.EnumMember); ok {
			if ok := enum.ValidMember(member); !ok {
				return exception.NewUmbraError("TY001", node, enum.Name.Lexeme, value)
			}
		}
	default:
		typeErr := types.CheckPrimitiveType(parsedType, value, nullable, node)

		if typeErr != nil {
			return typeErr
		}
	}

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
		} else {
			value = zero(stmt.Type.Type)
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

		return exception.NewUmbraError("RT039", stmt, types.SafeParseUmbraType(value))
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
		return Return{Value: value}
	case ast.FunctionExpression:
		parsedReturnType, parentEnum, err := parseRuntimeType(stmt.ReturnType, env)

		if err != nil {
			return err
		}

		env.Create(
			stmt,
			stmt.Name.Lexeme,
			FunctionDeclaration{Itself: &stmt, Environment: env, ReturnType: struct {
				Type   types.UmbraType
				Parent ast.EnumStatement
			}{Type: parsedReturnType, Parent: parentEnum}},
			types.FUN,
			false,
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

		stop, err := Evaluate(stmt.Stop, env)
		if err != nil {
			return err
		}

		var parsedStop float64
		var ok bool

		if parsedStop, ok = stop.(float64); !ok {
			return exception.NewUmbraError("RT022", stmt, types.SafeParseUmbraType(stop))
		}

		stepValue, err := Evaluate(stmt.Step, env)
		if err != nil {
			return err
		}

		step, ok := stepValue.(float64)
		if !ok {
			return exception.NewUmbraError("RT023", stmt, types.SafeParseUmbraType(stepValue))
		}

		for {
			loopEnv := environment.NewEnvironment(forEnv)
			controlVar, exists := loopEnv.Get(initializedVarName, true)
			if !exists {
				return exception.NewUmbraError("RT021", stmt, initializedVarName)
			}

			var condition bool
			if step >= 0 {
				condition = controlVar.Data.(float64) <= parsedStop
			} else {
				condition = controlVar.Data.(float64) >= parsedStop
			}

			if !condition {
				break
			}

			bodyErr := Interpret(stmt.Body, loopEnv)

			if _, ok := bodyErr.(Break); ok {
				break
			}

			loopEnv.Set(initializedVarName, controlVar.Data.(float64)+step)

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
				return exception.NewUmbraError("RT024", stmt, types.SafeParseUmbraType(parsedCondition))
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
				return exception.NewUmbraError("RT025", stmt, identifier.Lexeme)
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
			stmt,
			stmt.Name.Lexeme,
			stmt,
			types.ENUM,
			false,
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
				callback, err := processFunction(matchCase.Callback, env)

				if err != nil {
					return err
				}

				_, err = processFunctionCall(callback, value.(ast.EnumMember).Arguments, env)
				return err
			}
		}

		return nil
	default:
		return exception.NewUmbraError("RT000", stmt, reflect.TypeOf(statement).Name())
	}
}
