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

type FunctionStatement struct {
	Name   tokens.Token
	Params []tokens.Token
	Body   []Statement
}
type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

type PrintStatement struct {
	Expression Expression
}

type ReturnStatement struct {
	Keyword tokens.Token
	Value   Expression
}

type BreakStatement struct{}

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
	Name         tokens.Token
	Declarations []Statement
}
