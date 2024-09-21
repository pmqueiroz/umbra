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
	VisitBreakStatement(stmt BreakStatement)
	VisitVarStatement(stmt VarStatement)
	VisitPackageStatement(stmt ModuleStatement)
	VisitInitializedForStatement(stmt InitializedForStatement)
	VisitConditionalForStatement(stmt ConditionalForStatement)
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

type BreakStatement struct{}

func (stmt BreakStatement) Accept(visitor StatementVisitor) {
	visitor.VisitBreakStatement(stmt)
}

type VarStatement struct {
	Name        tokens.Token
	Initializer Expression
	Mutable     bool
	Type        tokens.Token
}

func (stmt VarStatement) Accept(visitor StatementVisitor) {
	visitor.VisitVarStatement(stmt)
}

type InitializedForStatement struct {
	Start Statement
	Stop  Expression
	Step  Expression
	Body  Statement
}

func (stmt InitializedForStatement) Accept(visitor StatementVisitor) {
	visitor.VisitInitializedForStatement(stmt)
}

type ConditionalForStatement struct {
	Condition Expression
	Body      Statement
}

func (stmt ConditionalForStatement) Accept(visitor StatementVisitor) {
	visitor.VisitConditionalForStatement(stmt)
}

type ModuleStatement struct {
	Name         tokens.Token
	Declarations []Statement
}

func (stmt ModuleStatement) Accept(visitor StatementVisitor) {
	visitor.VisitPackageStatement(stmt)
}
