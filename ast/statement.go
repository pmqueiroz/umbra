package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
)

type Statement interface {
}

type BlockStatement struct {
	Statements []Statement
}

type ExpressionStatement struct {
	Expression Expression
}

type Parameter struct {
	Name     tokens.Token
	Type     tokens.Token
	Variadic bool
}

type FunctionStatement struct {
	Name       tokens.Token
	Params     []Parameter
	ReturnType tokens.Token
	Body       []Statement
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

type PrintChannel int

const (
	StdoutChannel PrintChannel = iota
	StderrChannel
)

type PrintStatement struct {
	Expression Expression
	Channel    PrintChannel
}

type ReturnStatement struct {
	Keyword tokens.Token
	Value   Expression
}

type BreakStatement struct{}

type PublicStatement struct {
	Keyword    tokens.Token
	Identifier tokens.Token
}

type ImportStatement struct {
	Keyword tokens.Token
	Path    tokens.Token
}

type VarStatement struct {
	Name        tokens.Token
	Initializer Expression
	Mutable     bool
	Type        tokens.Token
}

type InitializedForStatement struct {
	Start Statement
	Stop  Expression
	Step  Expression
	Body  Statement
}

type ConditionalForStatement struct {
	Condition Expression
	Body      Statement
}

type ModuleStatement struct {
	Declarations []Statement
}
