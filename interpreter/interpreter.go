package interpreter

import (
	"fmt"
	"os"
	"strings"

	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/exception"
	"github.com/pmqueiroz/umbra/helpers"
	"github.com/pmqueiroz/umbra/tokens"
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
	Environment *Environment
}

func extractVarName(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case ast.VarStatement:
		return s.Name.Lexeme
	default:
		return ""
	}
}

func Interpret(statement ast.Statement, env *Environment) error {
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

		if str, ok := value.(string); ok {
			str = strings.ReplaceAll(str, "\\n", "\n")
			str = strings.ReplaceAll(str, "\\t", "\t")
			str = strings.ReplaceAll(str, "\\\"", "\"")
			str = strings.ReplaceAll(str, "\\\\", "\\")
			output = fmt.Sprint(str)
		} else {
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

			typeErr := CheckType(stmt.Type.Type, value, stmt.Nullable)

			if typeErr != nil {
				return typeErr
			}
		}

		env.Create(stmt.Name.Lexeme, value, stmt.Type.Type, stmt.Nullable)
		return nil
	case ast.BlockStatement:
		newEnv := NewEnvironment(env)
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
		env.Create(stmt.Name.Lexeme, FunctionDeclaration{Itself: &stmt, Environment: env}, tokens.FUN_TYPE, false)
		return nil
	case ast.ExpressionStatement:
		_, err := Evaluate(stmt.Expression, env)
		return err
	case ast.InitializedForStatement:
		forEnv := NewEnvironment(env)
		if err := Interpret(stmt.Start, forEnv); err != nil {
			return err
		}

		initializedVarName := extractVarName(stmt.Start)

		for {
			loopEnv := NewEnvironment(forEnv)
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
				condition = controlVar.data.(float64) <= parsedStop
			} else {
				return exception.NewRuntimeError("RT022", helpers.UmbraType(stop))
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
				return exception.NewRuntimeError("RT023", helpers.UmbraType(stepValue))
			}

			loopEnv.Set(initializedVarName, controlVar.data.(float64)+step)

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
			loopEnv := NewEnvironment(env)

			condition, err := Evaluate(stmt.Condition, loopEnv)
			if err != nil {
				return err
			}

			parsedCondition, ok := condition.(bool)
			if !ok {
				return exception.NewRuntimeError("RT024", helpers.UmbraType(parsedCondition))
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
		namespace, err := LoadModule(stmt.Path.Lexeme)

		if err != nil {
			return err
		}

		env.CreateNamespace(stmt.Path.Lexeme, &namespace)

		return nil
	default:
		return exception.NewRuntimeError("RT000", litter.Sdump(statement))
	}
}
