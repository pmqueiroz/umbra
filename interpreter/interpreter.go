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

type Break struct{}

func (r Break) Error() string {
	return "for loop break"
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
		if stmt.Initializer != nil {
			value, err = Evaluate(stmt.Initializer, env)
			if err != nil {
				return err
			}
		}
		env.Create(stmt.Name.Lexeme, value)
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
		env.Create(stmt.Name.Lexeme, FunctionDeclaration{Itself: &stmt, Environment: env})
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
			controlVar, ok := loopEnv.Get(initializedVarName)
			if !ok {
				return fmt.Errorf("control variable not found in environment: %s", initializedVarName)
			}

			stop, err := Evaluate(stmt.Stop, loopEnv)
			if err != nil {
				return err
			}

			var condition bool
			if parsedStop, ok := stop.(float64); ok {
				condition = controlVar.(float64) <= parsedStop
			} else {
				return fmt.Errorf("loop stop should be a number, got: %T", stop)
			}

			if !condition {
				break
			}

			if err := Interpret(stmt.Body, loopEnv); err != nil {
				if _, ok := err.(Break); ok {
					break
				}
				return err
			}

			stepValue, err := Evaluate(stmt.Step, loopEnv)
			if err != nil {
				return err
			}

			step, ok := stepValue.(float64)
			if !ok {
				return fmt.Errorf("loop step should be a number, got: %T", stepValue)
			}

			loopEnv.Set(initializedVarName, controlVar.(float64)+step)
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
				return fmt.Errorf("loop condition should be a boolean, got: %T", parsedCondition)
			}

			if !parsedCondition {
				break
			}

			if err := Interpret(stmt.Body, loopEnv); err != nil {
				if _, ok := err.(Break); ok {
					break
				}
				return err
			}
		}
		return nil
	case ast.BreakStatement:
		return Break{}
	default:
		return fmt.Errorf("unknown declaration: %T", statement)
	}
}
