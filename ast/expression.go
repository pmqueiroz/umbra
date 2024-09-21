package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
)

type Expression interface{}

type AssignExpression struct {
	Name  tokens.Token
	Value Expression
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
	Properties []KeyValueExpression
}

type KeyValueExpression struct {
	Key   tokens.Token
	Value Expression
}
