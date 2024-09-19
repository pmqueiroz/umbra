package ast

import (
	"fmt"

	"github.com/pmqueiroz/umbra/tokens"
)

type Parser struct {
	tokenList []tokens.Token
	current   int
}

func (p *Parser) peek() tokens.Token {
	return p.tokenList[p.current]
}

func (p *Parser) isAtEOF() bool {
	return p.peek().Id == tokens.EOF
}

func (p *Parser) previous() tokens.Token {
	return p.tokenList[p.current-1]
}

func (p *Parser) advance() tokens.Token {
	if !p.isAtEOF() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(tokenType tokens.TokenType) bool {
	if p.isAtEOF() {
		return false
	}

	return p.peek().Id == tokenType
}

func (p *Parser) match(types ...tokens.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType tokens.TokenType, message string) tokens.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(message)
}

func (p *Parser) block() Statement {
	var statements []Statement

	for !p.check(tokens.RIGHT_BRACE) && !p.isAtEOF() {
		statements = append(statements, p.declaration())
	}

	p.consume(tokens.RIGHT_BRACE, "Expect '}' after block.")

	return BlockStatement{
		Statements: statements,
	}
}

func (p *Parser) function() Statement {
	name := p.consume(tokens.IDENTIFIER, "Expect function name.")

	p.consume(tokens.LEFT_PAREN, "Expect '(' after function name.")

	var params []tokens.Token

	if !p.check(tokens.RIGHT_PAREN) {
		for {
			if len(params) >= 3 {
				panic("Can't have more than 3 parameters.")
			}

			params = append(params, p.consume(tokens.IDENTIFIER, "Expect parameter name."))

			if !p.match(tokens.COMMA) {
				break
			}
		}
	}

	p.consume(tokens.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(tokens.LEFT_BRACE, "Expect '{' before function body.")

	var body []Statement

	return FunctionStatement{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (p *Parser) primary() Expression {
	if p.match(tokens.FALSE) {
		return LiteralExpression{
			Value: "false",
		}
	}

	if p.match(tokens.TRUE) {
		return LiteralExpression{
			Value: "true",
		}
	}

	if p.match(tokens.NULL) {
		return LiteralExpression{
			Value: "null",
		}
	}

	if p.match(tokens.NUMERIC, tokens.STRING) {
		return LiteralExpression{
			Value: p.previous().Raw.Value,
		}
	}

	if p.match(tokens.LEFT_PAREN) {
		expr := p.assignment()
		p.consume(tokens.RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpression{
			Expression: expr,
		}
	}

	panic("Expect expression.")
}

func (p *Parser) finishCall(expr Expression) Expression {
	var arguments []Expression

	if !p.check(tokens.RIGHT_PAREN) {
		for {
			if len(arguments) >= 3 {
				panic("Can't have more than 3 arguments.")
			}

			arguments = append(arguments, p.assignment())

			if !p.match(tokens.COMMA) {
				break
			}
		}
	}

	paren := p.consume(tokens.RIGHT_PAREN, "Expect ')' after arguments.")

	return CallExpression{
		Callee:    expr,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) call() Expression {
	expr := p.primary()

	for {
		if p.match(tokens.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) unary() Expression {
	if p.match(tokens.NOT, tokens.MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpression{
			Operator: operator,
			Right:    right,
		}
	}

	return p.call()
}

func (p *Parser) multiplication() Expression {
	expr := p.unary()

	for p.match(tokens.SLASH, tokens.STAR) {
		expr = BinaryExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.unary(),
		}
	}

	return expr
}

func (p *Parser) addition() Expression {
	expr := p.multiplication()

	for p.match(tokens.MINUS, tokens.PLUS) {
		expr = BinaryExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.multiplication(),
		}
	}

	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.addition()

	for p.match(tokens.GREATER_THAN, tokens.GREATER_THAN_EQUAL, tokens.LESS_THAN, tokens.LESS_THAN_EQUAL) {
		expr = BinaryExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.addition(),
		}
	}

	return expr
}

func (p *Parser) equality() Expression {
	expr := p.comparison()

	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		expr = BinaryExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.comparison(),
		}
	}

	return expr
}

func (p *Parser) and() Expression {
	expr := p.equality()

	for p.match(tokens.AND) {
		expr = LogicalExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.equality(),
		}
	}

	return expr
}

func (p *Parser) or() Expression {
	expr := p.and()

	for p.match(tokens.OR) {
		expr = LogicalExpression{
			Left:     expr,
			Operator: p.previous(),
			Right:    p.and(),
		}
	}

	return expr
}

func (p *Parser) assignment() Expression {
	expr := p.or()

	if p.match(tokens.EQUAL) {
		value := p.assignment()

		if assign, ok := expr.(VariableExpression); ok {
			name := assign.Name
			return AssignExpression{
				Name:  name,
				Value: value,
			}
		}

		panic("Invalid assignment target.")
	}

	return expr
}

func (p *Parser) varDeclaration(tokenType tokens.TokenType) Statement {
	name := p.consume(tokens.IDENTIFIER, "Expect variable name.")

	var initializer Expression

	if p.match(tokens.EQUAL) {
		initializer = p.assignment()
	}

	return VarStatement{
		Name:        name,
		Initializer: initializer,
		Type:        tokenType,
	}
}

func (p *Parser) whileStatement() Statement {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.assignment()
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) returnStatement() Statement {
	keyword := p.previous()
	value := p.assignment()

	return ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) printStatement() Statement {
	value := p.assignment()
	return PrintStatement{
		Expression: value,
	}
}

func (p *Parser) ifStatement() Statement {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.assignment()
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Statement

	if p.match(tokens.ELSE) {
		elseBranch = p.statement()
	}

	return IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) forStatement() Statement {
	p.consume(tokens.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Statement
	if p.match(tokens.SEMICOLON) {
		initializer = nil
	} else if p.match(tokens.STR_VAR, tokens.NUM_VAR, tokens.OBJ_VAR, tokens.ARR_VAR) {
		initializer = p.varDeclaration(p.previous().Id)
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expression
	if !p.check(tokens.SEMICOLON) {
		condition = p.assignment()
	}
	p.consume(tokens.SEMICOLON, "Expect ';' after loop condition.")

	var increment Expression
	if !p.check(tokens.RIGHT_PAREN) {
		increment = p.assignment()
	}
	p.consume(tokens.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()
	if increment != nil {
		body = BlockStatement{
			Statements: []Statement{body, ExpressionStatement{Expression: increment}},
		}
	}

	if condition == nil {
		condition = LiteralExpression{Value: "true"}
	}
	body = WhileStatement{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = BlockStatement{
			Statements: []Statement{initializer, body},
		}
	}

	return body
}

func (p *Parser) expressionStatement() Statement {
	expr := p.assignment()
	return ExpressionStatement{
		Expression: expr,
	}
}

func (p *Parser) statement() Statement {
	if p.match(tokens.FOR) {
		return p.forStatement()
	}
	if p.match(tokens.IF) {
		return p.ifStatement()
	}
	if p.match(tokens.PRINT) {
		return p.printStatement()
	}
	if p.match(tokens.RETURN) {
		return p.returnStatement()
	}
	if p.match(tokens.WHILE) {
		return p.whileStatement()
	}
	if p.match(tokens.LEFT_BRACE) {
		return p.block()
	}

	return p.expressionStatement()
}

func (p *Parser) declaration() Statement {
	if p.match(tokens.FUN) {
		return p.function()
	}

	if p.match(tokens.STR_VAR) {
		return p.varDeclaration(tokens.STR_VAR)
	}

	if p.match(tokens.NUM_VAR) {
		return p.varDeclaration(tokens.NUM_VAR)
	}

	if p.match(tokens.OBJ_VAR) {
		return p.varDeclaration(tokens.OBJ_VAR)
	}

	if p.match(tokens.ARR_VAR) {
		return p.varDeclaration(tokens.ARR_VAR)
	}

	return p.statement()
}

func Parse(tokenList []tokens.Token) {
	var statements []Statement
	parser := Parser{
		tokenList: tokenList,
		current:   0,
	}

	for !parser.isAtEOF() {
		statements = append(statements, parser.declaration())
	}

	for _, stmt := range statements {
		fmt.Printf("%#v\n", stmt)
	}
}
