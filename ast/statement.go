package ast

import (
	"github.com/pmqueiroz/umbra/globals"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/pmqueiroz/umbra/types"
)

type Statement = globals.Node

type BlockStatement struct {
	Statements []Statement
}

func (s BlockStatement) Reference() string {
	statements := ""

	for _, statement := range s.Statements {
		statements += statement.Reference()
	}

	return "{" + statements + "}"
}

func (s BlockStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}
	for _, statement := range s.Statements {
		locs = append(locs, statement.GetLocs()...)
	}

	return locs
}

type ExpressionStatement struct {
	Expression Expression
}

func (s ExpressionStatement) Reference() string {
	return s.Expression.Reference()
}

func (s ExpressionStatement) GetLocs() []globals.Loc {
	return s.Expression.GetLocs()
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

func (s MatchStatement) Reference() string {
	return "match " + s.Expression.Reference() + " { ... }"
}

func (s MatchStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, s.Expression.GetLocs()...)

	return locs
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (s IfStatement) Reference() string {
	return "if " + s.Condition.Reference() + " { ... }"
}

func (s IfStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, s.Condition.GetLocs()...)

	return locs
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

func (s PrintStatement) Reference() string {
	channel := "stdout"

	if s.Channel == StderrChannel {
		channel = "stderr"
	}

	return channel + " " + s.Expression.Reference()
}

func (s PrintStatement) GetLocs() []globals.Loc {
	return s.Expression.GetLocs()
}

type ReturnStatement struct {
	Keyword tokens.Token
	Value   Expression
}

func (s ReturnStatement) Reference() string {
	return "return " + s.Value.Reference()
}

func (s ReturnStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{s.Keyword.Loc}

	locs = append(locs, s.Value.GetLocs()...)

	return locs
}

type BreakStatement struct {
	Loc globals.Loc
}

func (s BreakStatement) Reference() string {
	return "break"
}

func (s BreakStatement) GetLocs() []globals.Loc {
	return []globals.Loc{s.Loc}
}

type ContinueStatement struct {
	Loc globals.Loc
}

func (s ContinueStatement) Reference() string {
	return "continue"
}

func (s ContinueStatement) GetLocs() []globals.Loc {
	return []globals.Loc{s.Loc}
}

type PublicStatement struct {
	Keyword     tokens.Token
	Identifiers []tokens.Token
}

func (s PublicStatement) Reference() string {
	identifiers := ""

	for _, identifier := range s.Identifiers {
		identifiers += identifier.Lexeme
		identifiers += " "
	}

	return "pub {" + identifiers + "}"
}

func (s PublicStatement) GetLocs() []globals.Loc {
	return []globals.Loc{s.Keyword.Loc}
}

type ImportStatement struct {
	Keyword tokens.Token
	Path    tokens.Token
}

func (s ImportStatement) Reference() string {
	return "import " + s.Path.Lexeme
}

func (s ImportStatement) GetLocs() []globals.Loc {
	return []globals.Loc{s.Keyword.Loc}
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

func (s EnumStatement) Reference() string {
	return "enum " + s.Name.Lexeme + " { ... }"
}

func (s EnumStatement) GetLocs() []globals.Loc {
	return []globals.Loc{s.Name.Loc}
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

func (s VarStatement) Reference() string {
	varInit := "const"
	initializer := ""

	if s.Mutable {
		varInit = "mut"
	}

	if s.Initializer != nil {
		initializer = " = " + s.Initializer.Reference()
	}

	return varInit + " " + s.Name.Lexeme + " " + s.Type.Lexeme + initializer
}

func (s VarStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{s.Name.Loc}

	locs = append(locs, s.Initializer.GetLocs()...)
	locs = append(locs, s.Type.Loc)

	return locs
}

type ArrayDestructuringStatement struct {
	Declarations []VarStatement
	Expr         Expression
}

func (s ArrayDestructuringStatement) Reference() string {
	declarations := ""

	for index, declaration := range s.Declarations {
		declarations += declaration.Reference()

		if index < len(s.Declarations)-1 {
			declarations += ", "
		}
	}

	return declarations + " = " + s.Expr.Reference()
}

func (s ArrayDestructuringStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	for _, declaration := range s.Declarations {
		locs = append(locs, declaration.GetLocs()...)
	}

	return locs
}

type InitializedForStatement struct {
	Start Statement
	Stop  Expression
	Step  Expression
	Body  Statement
}

func (s InitializedForStatement) Reference() string {
	return "for " + s.Start.Reference() + ", " + s.Stop.Reference() + ", " + s.Step.Reference() + " { ... }"
}

func (s InitializedForStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, s.Start.GetLocs()...)
	locs = append(locs, s.Stop.GetLocs()...)
	locs = append(locs, s.Step.GetLocs()...)

	return locs
}

type ConditionalForStatement struct {
	Condition Expression
	Body      Statement
}

func (s ConditionalForStatement) Reference() string {
	return "for " + s.Condition.Reference() + " { ... }"
}

func (s ConditionalForStatement) GetLocs() []globals.Loc {
	locs := []globals.Loc{}

	locs = append(locs, s.Condition.GetLocs()...)

	return locs
}

type ModuleStatement struct {
	Declarations []Statement
}

func (s ModuleStatement) Reference() string {
	return "Module"
}

func (s ModuleStatement) GetLocs() []globals.Loc {
	return []globals.Loc{{}}
}
