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

func extractVarName(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case ast.VarStatement:
		return s.Name.Raw.Value
	default:
		return ""
	}
}

func Interpret(stmt ast.Statement, env *Environment) error {
	switch s := stmt.(type) {
	case ast.PrintStatement:
		value, err := Evaluate(s.Expression, env)
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
		if s.Initializer != nil {
			value, err = Evaluate(s.Initializer, env)
			if err != nil {
				return err
			}
		}
		env.Create(s.Name.Raw.Value, value)
		return nil
	case ast.BlockStatement:
		newEnv := NewEnvironment(env)
		for _, stmt := range s.Statements {
			if err := Interpret(stmt, newEnv); err != nil {
				return err
			}
		}
		return nil
	case ast.ModuleStatement:
		for _, stmt := range s.Declarations {
			if err := Interpret(stmt, env); err != nil {
				return err
			}
		}
		return nil
	case ast.IfStatement:
		condition, err := Evaluate(s.Condition, env)
		if err != nil {
			return err
		}

		if condition.(bool) {
			return Interpret(s.ThenBranch, env)
		} else if s.ElseBranch != nil {
			return Interpret(s.ElseBranch, env)
		}
		return nil
	case ast.ReturnStatement:
		value, err := Evaluate(s.Value, env)
		if err != nil {
			return err
		}
		return Return{value: value}
	case ast.FunctionStatement:
		env.Create(s.Name.Raw.Value, s)
		return nil
	case ast.ExpressionStatement:
		_, err := Evaluate(s.Expression, env)
		return err
	case ast.InitializedForStatement:
		forEnv := NewEnvironment(env)
		if err := Interpret(s.Start, forEnv); err != nil {
			return err
		}

		initializedVarName := extractVarName(s.Start)

		for {
			loopEnv := NewEnvironment(forEnv)
			controlVar, ok := loopEnv.Get(initializedVarName)
			if !ok {
				return fmt.Errorf("control variable not found in environment: %s", initializedVarName)
			}

			stop, err := Evaluate(s.Stop, loopEnv)
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

			if err := Interpret(s.Body, loopEnv); err != nil {
				if _, ok := err.(Break); ok {
					break
				}
				return err
			}

			stepValue, err := Evaluate(s.Step, loopEnv)
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

			condition, err := Evaluate(s.Condition, loopEnv)
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

			if err := Interpret(s.Body, loopEnv); err != nil {
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
		return fmt.Errorf("unknown declaration: %T", stmt)
	}
}
