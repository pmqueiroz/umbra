package ast

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pmqueiroz/umbra/exception"
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
	return p.peek().Type == tokens.EOF
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

	return p.peek().Type == tokenType
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

	p.throw(errorMessage)
	return tokens.Token{}
}

func (p *Parser) throw(message string) {
	current_token := p.peek()
	fmt.Println(exception.NewSyntaxError(
		message,
		current_token.Loc.Line,
		current_token.Loc.Range.From,
		fmt.Sprintf("%#v", p.peek()),
	))
	os.Exit(1)
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

	var params []Parameter

	if !p.check(tokens.RIGHT_PARENTHESIS) {
		for {
			paramName := p.consume("Expect parameter name.", tokens.IDENTIFIER)
			paramType := p.consume("Expect parameter type.", tokens.DATA_TYPES...)
			params = append(params, Parameter{
				Name: paramName,
				Type: paramType,
			})

			if !p.match(tokens.COMMA) {
				break
			}
		}
	}

	p.consume("Expect ')' after parameters.", tokens.RIGHT_PARENTHESIS)

	var returnType tokens.Token

	if !p.check(tokens.LEFT_BRACE) {
		returnType = p.consume("Expect return type.", append([]tokens.TokenType{tokens.VOID_TYPE}, tokens.DATA_TYPES...)...)
	} else {
		currentToken := p.peek()
		returnType = tokens.Token{
			Type:   tokens.VOID_TYPE,
			Loc:    currentToken.Loc,
			Lexeme: currentToken.Lexeme,
		}
	}

	p.consume("Expect '{' before function body.", tokens.LEFT_BRACE)

	_, body := p.block()

	return FunctionStatement{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *Parser) array() Expression {
	var elements []Expression

	if !p.check(tokens.RIGHT_BRACKET) {
		for {
			elements = append(elements, p.expression())

			if !p.match(tokens.COMMA) || p.check(tokens.RIGHT_BRACKET) {
				break
			}
		}
	}

	p.consume("Expect ']' after expression.", tokens.RIGHT_BRACKET)

	return ArrayExpression{
		Elements: elements,
	}
}

func (p *Parser) hashmap() Expression {
	properties := make(map[Expression]Expression)

	if !p.check(tokens.RIGHT_BRACE) {
		for {
			name := p.consume("Expect property name.", tokens.IDENTIFIER)
			p.consume("Expect ':' after property identifier in hashmap", tokens.COLON)

			properties[LiteralExpression{Value: name.Lexeme}] = p.expression()

			if !p.match(tokens.COMMA) || p.check(tokens.RIGHT_BRACE) {
				break
			}
		}
	}

	p.consume("Expect '}' after expression.", tokens.RIGHT_BRACE)

	return HashmapExpression{
		Pairs: properties,
	}
}

func (p *Parser) numeric() Expression {
	value, err := strconv.ParseFloat(p.previous().Lexeme, 64)

	if err != nil {
		p.throw("Unable to convert number.")
	}

	return LiteralExpression{
		Value: value,
	}
}

func (p *Parser) primary() Expression {
	if p.match(tokens.FALSE) {
		return LiteralExpression{
			Value: false,
		}
	}

	if p.match(tokens.TRUE) {
		return LiteralExpression{
			Value: true,
		}
	}

	if p.match(tokens.NULL) {
		return LiteralExpression{
			Value: nil,
		}
	}

	if p.match(tokens.NUMERIC) {
		return p.numeric()
	}

	if p.match(tokens.STRING) {
		return LiteralExpression{
			Value: p.previous().Lexeme,
		}
	}

	if p.match(tokens.IDENTIFIER) {
		return VariableExpression{
			Name: p.previous(),
		}
	}

	if p.match(tokens.LEFT_BRACE) {
		return p.hashmap()
	}

	if p.match(tokens.LEFT_BRACKET) {
		return p.array()
	}

	if p.match(tokens.LEFT_PARENTHESIS) {
		expr := p.expression()
		p.consume("Expect ')' after expression.", tokens.RIGHT_PARENTHESIS)
		return GroupingExpression{
			Expression: expr,
		}
	}

	p.throw("Expect expression.")
	return LiteralExpression{
		Value: nil,
	}
}

func (p *Parser) finishCall(expr Expression) Expression {
	var arguments []Expression

	if !p.check(tokens.RIGHT_PARENTHESIS) {
		for {
			if p.check(tokens.RIGHT_PARENTHESIS) {
				break
			}

			arguments = append(arguments, p.expression())

			if !p.match(tokens.COMMA) {
				break
			}
		}
	}

	p.consume("Expect ')' after arguments.", tokens.RIGHT_PARENTHESIS)

	return CallExpression{
		Callee:    expr,
		Arguments: arguments,
	}
}

func (p *Parser) call() Expression {
	expr := p.primary()

	for {
		if p.match(tokens.LEFT_PARENTHESIS) {
			return p.finishCall(expr)
		} else if p.match(tokens.DOT) {
			property := p.consume("Expect property name after '.'.", tokens.IDENTIFIER)
			expr = MemberExpression{
				Object:   expr,
				Property: property,
			}
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

func (p *Parser) expression() Expression {
	expr := p.or()

	if p.match(tokens.EQUAL) {
		value := p.expression()

		switch target := expr.(type) {
		case VariableExpression, MemberExpression:
			return AssignExpression{
				Target: target,
				Value:  value,
			}
		default:
			p.throw("Invalid assignment target.")
		}

	}

	return expr
}

func (p *Parser) varDeclaration() Statement {
	isMutable := p.previous().Type == tokens.MUT
	name := p.consume("Expect variable name.", tokens.IDENTIFIER)
	variableType := p.consume("Expect variable type.", tokens.DATA_TYPES...)

	var initializer Expression

	if p.match(tokens.EQUAL) {
		initializer = p.expression()
	}

	return VarStatement{
		Name:        name,
		Initializer: initializer,
		Mutable:     isMutable,
		Type:        variableType,
	}
}

func (p *Parser) breakStatement() Statement {
	return BreakStatement{}
}

func (p *Parser) returnStatement() Statement {
	keyword := p.previous()
	value := p.expression()

	return ReturnStatement{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) publicStatement() Statement {
	keyword := p.previous()
	name := p.consume("Expect variable name.", tokens.IDENTIFIER)

	return PublicStatement{
		Keyword:    keyword,
		Identifier: name,
	}
}

func (p *Parser) printStatement(channel PrintChannel) Statement {
	value := p.expression()
	return PrintStatement{
		Expression: value,
		Channel:    channel,
	}
}

func (p *Parser) ifStatement() Statement {
	condition := p.expression()

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

func (p *Parser) initializedForStatement() Statement {
	start := p.varDeclaration()
	var step Expression

	p.consume("Expect ',' after start.", tokens.COMMA)

	stop := p.expression()

	switch stop.(type) {
	case LiteralExpression, VariableExpression:
	default:
		p.throw("Expect stop to be a numeric literal or variable declaration.")
	}

	if p.match(tokens.COMMA) {
		p.consume("Expect step after stop.", tokens.NUMERIC)
		step = p.numeric()
	} else {
		step = LiteralExpression{
			Value: 1.0,
		}
	}

	body := p.statement()

	return InitializedForStatement{
		Start: start,
		Stop:  stop,
		Step:  step,
		Body:  body,
	}

}

func (p *Parser) conditionalForStatement() Statement {
	var expr Expression

	if p.check(tokens.LEFT_BRACE) {
		expr = LiteralExpression{
			Value: true,
		}
	} else {
		expr = p.expression()
	}

	body := p.statement()

	return ConditionalForStatement{
		Condition: expr,
		Body:      body,
	}
}

func (p *Parser) forStatement() Statement {
	if p.match(tokens.CONST, tokens.MUT) {
		return p.initializedForStatement()
	}

	return p.conditionalForStatement()
}

func (p *Parser) expressionStatement() Statement {
	expr := p.expression()
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
	if p.match(tokens.STDOUT, tokens.STDERR) {
		channel := StdoutChannel
		if p.previous().Type == tokens.STDERR {
			channel = StderrChannel
		}
		return p.printStatement(channel)
	}
	if p.match(tokens.RETURN) {
		return p.returnStatement()
	}
	if p.match(tokens.BREAK) {
		return p.breakStatement()
	}
	if p.match(tokens.PUBLIC) {
		return p.publicStatement()
	}
	if p.match(tokens.LEFT_BRACE) {
		blockStatement, _ := p.block()
		return blockStatement
	}

	return p.expressionStatement()
}

func (p *Parser) declaration() Statement {
	if p.match(tokens.FUN) {
		return p.function()
	}

	if p.match(tokens.CONST, tokens.MUT) {
		return p.varDeclaration()
	}

	return p.statement()
}

func Parse(tokenList []tokens.Token) ModuleStatement {
	var declarations []Statement
	parser := Parser{
		tokenList: tokenList,
		current:   0,
	}

	for !parser.isAtEOF() {
		declarations = append(declarations, parser.declaration())
	}

	module := ModuleStatement{
		Declarations: declarations,
	}

	return module
}
