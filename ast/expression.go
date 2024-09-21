package ast

import (
	"github.com/pmqueiroz/umbra/tokens"
)

type ExpressionVisitor interface {
	VisitAssignExpression(expr AssignExpression)
	VisitBinaryExpression(expr Expression)
	VisitCallExpression(expr CallExpression)
	VisitGroupingExpression(expr GroupingExpression)
	VisitLiteralExpression(expr LiteralExpression)
	VisitHashmapExpression(expr HashmapExpression)
	VisitArrayExpression(expr ArrayExpression)
	VisitLogicalExpression(expr LogicalExpression)
	VisitUnaryExpression(expr UnaryExpression)
	VisitVariableExpression(expr VariableExpression)
	VisitKeyValueExpression(expr KeyValueExpression)
}

type Expression interface {
	Accept(visitor ExpressionVisitor)
}

type AssignExpression struct {
	Name  tokens.Token
	Value Expression
}

func (expr AssignExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitAssignExpression(expr)
}

type BinaryExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (expr BinaryExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitBinaryExpression(expr)
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
}

func (expr CallExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitCallExpression(expr)
}

type GroupingExpression struct {
	Expression Expression
}

func (expr GroupingExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitGroupingExpression(expr)
}

type StringLiteral struct {
	Value string
}

type LiteralExpression struct {
	Value interface{}
}

func (expr LiteralExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitLiteralExpression(expr)
}

type LogicalExpression struct {
	Left     Expression
	Operator tokens.Token
	Right    Expression
}

func (expr LogicalExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitLogicalExpression(expr)
}

type UnaryExpression struct {
	Operator tokens.Token
	Right    Expression
}

func (expr UnaryExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitUnaryExpression(expr)
}

type VariableExpression struct {
	Name tokens.Token
}

func (expr VariableExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitVariableExpression(expr)
}

type ArrayExpression struct {
	Elements []Expression
}

func (expr ArrayExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitArrayExpression(expr)
}

type HashmapExpression struct {
	Properties []KeyValueExpression
}

func (expr HashmapExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitHashmapExpression(expr)
}

type KeyValueExpression struct {
	Key   tokens.Token
	Value Expression
}

func (expr KeyValueExpression) Accept(visitor ExpressionVisitor) {
	visitor.VisitKeyValueExpression(expr)
}
