package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
)

type StatementVisitor interface {
	VisitBlockStatement(stmt BlockStatement)
	VisitExpressionStatement(stmt ExpressionStatement)
	VisitFunctionStatement(stmt FunctionStatement)
	VisitIfStatement(stmt IfStatement)
	VisitPrintStatement(stmt PrintStatement)
	VisitReturnStatement(stmt ReturnStatement)
	VisitVarStatement(stmt VarStatement)
	VisitWhileStatement(stmt WhileStatement)
	VisitPackageStatement(stmt PackageStatement)
}

type Statement interface {
	Accept(visitor StatementVisitor)
}

type BlockStatement struct {
	Statements []Statement
}

func (stmt BlockStatement) Accept(visitor StatementVisitor) {
	visitor.VisitBlockStatement(stmt)
}

type ExpressionStatement struct {
	Expression Expression
}

func (stmt ExpressionStatement) Accept(visitor StatementVisitor) {
	visitor.VisitExpressionStatement(stmt)
}

type FunctionStatement struct {
	Name   tokens.Token
	Params []tokens.Token
	Body   []Statement
}

func (stmt FunctionStatement) Accept(visitor StatementVisitor) {
	visitor.VisitFunctionStatement(stmt)
}

type IfStatement struct {
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (stmt IfStatement) Accept(visitor StatementVisitor) {
	visitor.VisitIfStatement(stmt)
}

type PrintStatement struct {
	Expression Expression
}

func (stmt PrintStatement) Accept(visitor StatementVisitor) {
	visitor.VisitPrintStatement(stmt)
}

type ReturnStatement struct {
	Keyword tokens.Token
	Value   Expression
}

func (stmt ReturnStatement) Accept(visitor StatementVisitor) {
	visitor.VisitReturnStatement(stmt)
}

type VarStatement struct {
	Name        tokens.Token
	Initializer Expression
	Type        tokens.TokenType
}

func (stmt VarStatement) Accept(visitor StatementVisitor) {
	visitor.VisitVarStatement(stmt)
}

type WhileStatement struct {
	Condition Expression
	Body      Statement
}

func (stmt WhileStatement) Accept(visitor StatementVisitor) {
	visitor.VisitWhileStatement(stmt)
}

type PackageStatement struct {
	Name tokens.Token
}

func (stmt PackageStatement) Accept(visitor StatementVisitor) {
	visitor.VisitPackageStatement(stmt)
}
