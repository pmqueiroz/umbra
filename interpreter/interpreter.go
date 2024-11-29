package interpreter

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
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
}

func extractVarName(stmt ast.Statement) string {
	switch s := stmt.(type) {
	case ast.VarStatement:
		return s.Name.Lexeme
	default:
		return ""
	}
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

			parsedType, enum, parseTypeErr := func() (types.UmbraType, ast.EnumStatement, error) {
				switch stmt.Type.Type {
				case tokens.IDENTIFIER:
					value, ok := env.Get(stmt.Type.Lexeme, true)
					if !ok {
						return types.UNKNOWN, ast.EnumStatement{}, exception.NewRuntimeError("RT002", stmt.Type.Lexeme)
					}

					if enum, ok := value.Data.(ast.EnumStatement); ok {
						return types.ENUM, enum, nil
					}

					return types.UNKNOWN, ast.EnumStatement{}, exception.NewRuntimeError("RT035", stmt.Type.Lexeme)
				default:
					t, e := types.ParseTokenType(stmt.Type.Type)
					return t, ast.EnumStatement{}, e
				}
			}()

			if parseTypeErr != nil {
				return parseTypeErr
			}

			switch parsedType {
			case types.ENUM:
				if member, ok := value.(ast.EnumMember); ok {
					if ok := enum.ValidMember(member); ok {
						return nil
					}
				}

				return exception.NewTypeError(fmt.Sprintf("expected %s got %s", enum.Name.Lexeme, value))
			default:
				typeErr := types.CheckPrimitiveType(parsedType, value, stmt.Nullable)

				if typeErr != nil {
					return typeErr
				}
			}

		}

		env.Create(stmt.Name.Lexeme, value, types.SafeParseUmbraType(stmt.Type.Type), stmt.Nullable, false)
		return nil
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
		env.Create(
			stmt.Name.Lexeme,
			FunctionDeclaration{Itself: &stmt, Environment: env},
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
	default:
		return exception.NewRuntimeError("RT000", litter.Sdump(statement))
	}
}
