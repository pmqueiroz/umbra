package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
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
	Type     types.UmbraType
	Variadic bool
	Nullable bool
}

type MatchCaseParameter struct {
	Name tokens.Token
}

type MatchCase struct {
	Expression Expression
	Callback   FunctionExpression
}

type MatchStatement struct {
	Expression Expression
	Cases      []MatchCase
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

type ContinueStatement struct{}

type PublicStatement struct {
	Keyword     tokens.Token
	Identifiers []tokens.Token
}

type ImportStatement struct {
	Keyword tokens.Token
	Path    tokens.Token
}

type EnumArgument struct {
	Type  types.UmbraType
	Value interface{}
}

type EnumMember struct {
	Name      string
	Arguments []EnumArgument
	Signature string
}

type EnumStatement struct {
	Name      tokens.Token
	Members   map[string]EnumMember
	Signature string
}

func (e *EnumStatement) Get(name tokens.Token) (EnumMember, bool) {
	member, ok := e.Members[name.Lexeme]

	return member, ok
}

func (e *EnumStatement) ValidMember(member EnumMember) bool {
	return e.Signature == member.Signature
}

type VarStatement struct {
	Name        tokens.Token
	Initializer Expression
	Mutable     bool
	Type        tokens.Token
	Nullable    bool
}

type ArrayDestructuringStatement struct {
	Declarations []VarStatement
	Expr         Expression
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
