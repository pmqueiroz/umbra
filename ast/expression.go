package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
)

type Expression interface{}

type AssignExpression struct {
	Target Expression
	Value  Expression
}

type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

type GroupingExpression struct {
	Expression Expression
}

type LiteralExpression struct {
	Value interface{}
}

type TypeConversionExpression struct {
	Type  tokens.Token
	Value Expression
}

type NaNExpression struct{}

type LogicalExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

type UnaryExpression struct {
	Operator tokens.Token
	Right    Expression
}

type VariableExpression struct {
	Name tokens.Token
}

type ArrayExpression struct {
	Elements []Expression
}

type HashmapExpression struct {
	Pairs map[Expression]Expression
}

type MemberExpression struct {
	Object   Expression
	Property Expression
	Computed bool
}

type NamespaceMemberExpression struct {
	Namespace Expression
	Property  tokens.Token
}

type SizeExpression struct {
	Value Expression
}

type FunctionExpression struct {
	Name       tokens.Token
	Params     []Parameter
	ReturnType tokens.Token
	Body       []Statement
}
