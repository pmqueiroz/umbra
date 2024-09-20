package ast

import (
	"fmt"

	umbra_error "github.com/pmqueiroz/umbra/error"
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

func (p *Parser) consume(errorMessage string, types ...tokens.TokenType) tokens.Token {
	for _, t := range types {
		if p.check(t) {
			return p.advance()
		}
	}

	panic(errorMessage)
}

func (p *Parser) block() (Statement, []Statement) {
	var statements []Statement

	for !p.check(tokens.RIGHT_BRACE) && !p.isAtEOF() {
		statements = append(statements, p.declaration())
	}

	p.consume("Expect '}' after block.", tokens.RIGHT_BRACE)

	return BlockStatement{
		Statements: statements,
	}, statements
}

func (p *Parser) function() Statement {
	name := p.consume("Expect function name.", tokens.IDENTIFIER)

	p.consume("Expect '(' after function name.", tokens.LEFT_PARENTHESIS)

	var params []tokens.Token

	if !p.check(tokens.RIGHT_PARENTHESIS) {
		for {
			if len(params) >= 3 {
				panic("Can't have more than 3 parameters.")
			}

			params = append(params, p.consume("Expect parameter name.", tokens.IDENTIFIER))

			if !p.match(tokens.COMMA) {
				break
			}
		}
	}

	p.consume("Expect ')' after parameters.", tokens.RIGHT_PARENTHESIS)

	p.consume("Expect '{' before function body.", tokens.LEFT_BRACE)

	_, body := p.block()

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

	if p.match(tokens.LEFT_PARENTHESIS) {
		expr := p.assignment()
		p.consume("Expect ')' after expression.", tokens.RIGHT_PARENTHESIS)
		return GroupingExpression{
			Expression: expr,
		}
	}

	current_token := p.peek()

	panic(
		umbra_error.NewSyntaxError("Expect expression.", current_token.Raw.Line, current_token.Raw.Column, fmt.Sprintf("%#v", p.peek())),
	)
}

func (p *Parser) finishCall(expr Expression) Expression {
	var arguments []Expression

	if !p.check(tokens.RIGHT_PARENTHESIS) {
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

	parenthesis := p.consume("Expect ')' after arguments.", tokens.RIGHT_PARENTHESIS)

	return CallExpression{
		Callee:      expr,
		Parenthesis: parenthesis,
		Arguments:   arguments,
	}
}

func (p *Parser) call() Expression {
	expr := p.primary()

	for {
		if p.match(tokens.LEFT_PARENTHESIS) {
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

func (p *Parser) packageDeclaration() Statement {
	name := p.consume("Expect package name.", tokens.IDENTIFIER)

	return PackageStatement{
		Name: name,
	}
}

func (p *Parser) varDeclaration() Statement {
	isMutable := p.previous().Id == tokens.MUT
	name := p.consume("Expect variable name.", tokens.IDENTIFIER)
	variableType := p.consume("Expect variable type.", tokens.STR_TYPE, tokens.NUM_TYPE, tokens.HASHMAP_TYPE, tokens.ARR_TYPE)

	var initializer Expression

	if p.match(tokens.EQUAL) {
		initializer = p.assignment()
	}

	return VarStatement{
		Name:        name,
		Initializer: initializer,
		Mutable:     isMutable,
		Type:        variableType,
	}
}

func (p *Parser) whileStatement() Statement {
	p.consume("Expect '(' after 'while'.", tokens.LEFT_PARENTHESIS)
	condition := p.assignment()
	p.consume("Expect ')' after condition.", tokens.RIGHT_PARENTHESIS)
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
	p.consume("Expect '(' after 'if'.", tokens.LEFT_PARENTHESIS)
	condition := p.assignment()
	p.consume("Expect ')' after if condition.", tokens.RIGHT_PARENTHESIS)

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
	p.consume("Expect '(' after 'for'.", tokens.LEFT_PARENTHESIS)

	var initializer Statement
	if p.match(tokens.SEMICOLON) {
		initializer = nil
	} else if p.match(tokens.STR_TYPE, tokens.NUM_TYPE, tokens.HASHMAP_TYPE, tokens.ARR_TYPE) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expression
	if !p.check(tokens.SEMICOLON) {
		condition = p.assignment()
	}
	p.consume("Expect ';' after loop condition.", tokens.SEMICOLON)

	var increment Expression
	if !p.check(tokens.RIGHT_PARENTHESIS) {
		increment = p.assignment()
	}
	p.consume("Expect ')' after for clauses.", tokens.RIGHT_PARENTHESIS)

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
		blockStatement, _ := p.block()
		return blockStatement
	}

	return p.expressionStatement()
}

func (p *Parser) declaration() Statement {
	if p.match(tokens.PACKAGE) {
		return p.packageDeclaration()
	}

	if p.match(tokens.FUN) {
		return p.function()
	}

	if p.match(tokens.VAR, tokens.MUT) {
		return p.varDeclaration()
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
