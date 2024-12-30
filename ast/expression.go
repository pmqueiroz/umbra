package ast

import (
	"fmt"

	"github.com/pmqueiroz/umbra/globals"
	"github.com/pmqueiroz/umbra/tokens"
)

type Expression = globals.Node

type AssignExpression struct {
	Target Expression
	Value  Expression
}

func (e AssignExpression) Reference() string {
	return e.Target.Reference() + " = " + e.Value.Reference()
}

type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e BinaryExpression) Reference() string {
	return fmt.Sprintf("%s %s %s", e.Left.Reference(), e.Operator.Lexeme, e.Right.Reference())
}

type IsExpression struct {
	Expr     Expression
	Expected tokens.Token
}

func (e IsExpression) Reference() string {
	return e.Expr.Reference() + " is " + e.Expected.Lexeme
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (e CallExpression) Reference() string {
	arguments := ""
	for i, element := range e.Arguments {
		arguments += element.Reference()
		if i < len(e.Arguments)-1 {
			arguments += ", "
		}
	}

	return e.Callee.Reference() + "(" + arguments + ")"
}

type GroupingExpression struct {
	Expression Expression
}

func (e GroupingExpression) Reference() string {
	return "(" + e.Expression.Reference() + ")"
}

type LiteralExpression struct {
	Lexeme string
	Value  interface{}
}

func (e LiteralExpression) Reference() string {
	return e.Lexeme
}

type TypeConversionExpression struct {
	Type  tokens.Token
	Value Expression
}

func (e TypeConversionExpression) Reference() string {
	return e.Type.Lexeme + "(" + e.Value.Reference() + ")"
}

type LogicalExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e LogicalExpression) Reference() string {
	return e.Left.Reference() + " " + e.Operator.Lexeme + " " + e.Right.Reference()
}

type UnaryExpression struct {
	Operator tokens.Token
	Right    Expression
}

func (e UnaryExpression) Reference() string {
	return e.Operator.Lexeme + e.Right.Reference()
}

type VariableExpression struct {
	Name tokens.Token
}

func (e VariableExpression) Reference() string {
	return e.Name.Lexeme
}

type ArrayExpression struct {
	Elements []Expression
}

func (e ArrayExpression) Reference() string {
	elements := ""
	for i, element := range e.Elements {
		elements += element.Reference()
		if i < len(e.Elements)-1 {
			elements += ", "
		}
	}

	return "[" + elements + "]"
}

type HashmapExpression struct {
	Pairs map[Expression]Expression
}

func (e HashmapExpression) Reference() string {
	arguments := ""
	index := 0
	for key, value := range e.Pairs {
		arguments += key.Reference() + ": " + value.Reference()
		if index < len(e.Pairs)-1 {
			arguments += ", "
		}

		index++
	}
	return "{" + arguments + "}"
}

type MemberExpressionType string

const (
	DotMember     MemberExpressionType = "DotMember"
	BracketMember MemberExpressionType = "BracketMember"
)

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
	Type     MemberExpressionType
}

func (e MemberExpression) Reference() string {
	if e.Type == DotMember {
		return e.Object.Reference() + "." + e.Property.Reference()
	}

	return e.Object.Reference() + "[" + e.Property.Reference() + "]"
}

type NamespaceMemberExpression struct {
	Namespace Expression
	Property  tokens.Token
}

func (e NamespaceMemberExpression) Reference() string {
	return e.Namespace.Reference() + "::" + e.Property.Lexeme
}

type FunctionExpression struct {
	Name       tokens.Token
	Params     []Parameter
	ReturnType tokens.Token
	Body       []Statement
}

func (e FunctionExpression) Reference() string {
	params := ""

	for i, param := range e.Params {
		if param.Variadic {
			params += "..."
		}
		params += param.Name.Lexeme + " " + string(param.Type)
		if i < len(e.Params)-1 {
			params += ", "
		}
	}

	return "def " + e.Name.Lexeme + "(" + params + ") " + e.ReturnType.Lexeme + " { ... }"
}
