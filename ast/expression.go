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

func (e AssignExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Target.GetLocs()...)
	locs = append(locs, e.Value.GetLocs()...)
	return locs
}

type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e BinaryExpression) Reference() string {
	return fmt.Sprintf("%s %s %s", e.Left.Reference(), e.Operator.Lexeme, e.Right.Reference())
}

func (e BinaryExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Left.GetLocs()...)
	locs = append(locs, e.Operator.Loc)
	locs = append(locs, e.Right.GetLocs()...)

	return locs
}

type IsExpression struct {
	Expr     Expression
	Expected tokens.Token
}

func (e IsExpression) Reference() string {
	return e.Expr.Reference() + " is " + e.Expected.Lexeme
}

func (e IsExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Expr.GetLocs()...)
	locs = append(locs, e.Expr.GetLocs()...)
	return locs
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

func (e CallExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Callee.GetLocs()...)
	arguments := []globals.Loc{}

	for _, argument := range e.Arguments {
		arguments = append(arguments, argument.GetLocs()...)
	}

	locs = append(locs, arguments...)
	return locs
}

type GroupingExpression struct {
	Expression Expression
}

func (e GroupingExpression) Reference() string {
	return "(" + e.Expression.Reference() + ")"
}

func (e GroupingExpression) GetLocs() []globals.Loc {
	return e.Expression.GetLocs()
}

type LiteralExpression struct {
	Loc    globals.Loc
	Lexeme string
	Value  interface{}
}

func (e LiteralExpression) Reference() string {
	return e.Lexeme
}

func (e LiteralExpression) GetLocs() []globals.Loc {
	return []globals.Loc{e.Loc}
}

type TypeConversionExpression struct {
	Type  tokens.Token
	Value Expression
}

func (e TypeConversionExpression) Reference() string {
	return e.Type.Lexeme + "(" + e.Value.Reference() + ")"
}

func (e TypeConversionExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Type.Loc)
	locs = append(locs, e.Value.GetLocs()...)
	return locs
}

type LogicalExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (e LogicalExpression) Reference() string {
	return e.Left.Reference() + " " + e.Operator.Lexeme + " " + e.Right.Reference()
}

func (e LogicalExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Left.GetLocs()...)
	locs = append(locs, e.Operator.Loc)
	locs = append(locs, e.Right.GetLocs()...)

	return locs
}

type UnaryExpression struct {
	Operator tokens.Token
	Right    Expression
}

func (e UnaryExpression) Reference() string {
	return e.Operator.Lexeme + e.Right.Reference()
}

func (e UnaryExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, e.Operator.Loc)
	locs = append(locs, e.Right.GetLocs()...)

	return locs
}

type VariableExpression struct {
	Name tokens.Token
}

func (e VariableExpression) Reference() string {
	return e.Name.Lexeme
}

func (e VariableExpression) GetLocs() []globals.Loc {
	return []globals.Loc{e.Name.Loc}
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

func (e ArrayExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	for _, element := range e.Elements {
		locs = append(locs, element.GetLocs()...)
	}
	return locs
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

func (e HashmapExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	for key, value := range e.Pairs {
		locs = append(locs, key.GetLocs()...)
		locs = append(locs, value.GetLocs()...)
	}
	return locs
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

func (e MemberExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}
	locs = append(locs, e.Object.GetLocs()...)
	locs = append(locs, e.Property.GetLocs()...)

	return locs
}

type NamespaceMemberExpression struct {
	Namespace Expression
	Property  tokens.Token
}

func (e NamespaceMemberExpression) Reference() string {
	return e.Namespace.Reference() + "::" + e.Property.Lexeme
}

func (e NamespaceMemberExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}
	locs = append(locs, e.Namespace.GetLocs()...)
	locs = append(locs, e.Property.Loc)

	return locs
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

func (e FunctionExpression) GetLocs() []globals.Loc {
	locs := []globals.Loc{}
	locs = append(locs, e.Name.Loc)
	paramsLocs := []globals.Loc{}
	for _, param := range e.Params {
		paramsLocs = append(paramsLocs, param.Name.Loc)
	}
	locs = append(locs, paramsLocs...)
	locs = append(locs, e.ReturnType.Loc)

	return locs
}
