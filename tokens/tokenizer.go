package tokens

import (
	"fmt"
	"unicode"

	"github.com/pmqueiroz/umbra/exception"
)

func isAlpha(char rune) bool {
	return unicode.IsLetter(char) || char == '_'
}

func isDigit(char rune) bool {
	return unicode.IsDigit(char)
}

func isAlphaNumeric(char rune) bool {
	return isAlpha(char) || unicode.IsDigit(char)
}

type Tokenizer struct {
	tokens                               []Token
	current, beginOfLexeme, line, column int
	source                               string
}

func (t *Tokenizer) isAtEnd() bool {
	return t.current >= len(t.source)
}

func (t *Tokenizer) add(token Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokenizer) advance() rune {
	t.current++
	t.column++
	return []rune(t.source)[t.current-1]
}

func (t *Tokenizer) previous() rune {
	return []rune(t.source)[t.current-1]
}

func (t *Tokenizer) peek() rune {
	if t.isAtEnd() {
		return '\000'
	}

	return []rune(t.source)[t.current]
}

func (t *Tokenizer) peekNext() rune {
	if t.current+1 >= len(t.source) {
		return '\000'
	}

	return []rune(t.source)[t.current+1]
}

func (t *Tokenizer) match(expected rune) bool {
	if t.isAtEnd() {
		return false
	}

	if []rune(t.source)[t.current] != expected {
		return false
	}

	t.current++
	return true
}

func (t *Tokenizer) addNonLiteralToken(tokenType TokenType) {
	t.add(Token{
		Id: tokenType,
		Raw: RawToken{
			Value:  string([]rune(t.source)[t.beginOfLexeme:t.current]),
			Line:   t.line,
			Column: t.column,
		},
	})
}

func (t *Tokenizer) advanceLine() {
	t.line++
	t.column = 0
}

func (t *Tokenizer) string() {
	for (t.peek() != '"' || (t.peek() == '"' && t.previous() == '\\')) && !t.isAtEnd() {
		if t.peek() == '\n' {
			t.advanceLine()
		}

		t.advance()
	}

	if t.isAtEnd() {
		fmt.Println(
			exception.NewSyntaxError("Unterminated string", t.line, t.column, string([]rune(t.source)[t.beginOfLexeme:t.current])),
		)
		return
	}

	t.advance()

	t.add(Token{
		Id: STRING,
		Raw: RawToken{
			Value:  string([]rune(t.source)[t.beginOfLexeme+1 : t.current-1]),
			Line:   t.line,
			Column: t.column,
		},
	})
}

func (t *Tokenizer) numeric() {
	for isDigit(t.peek()) {
		t.advance()
	}

	if t.peek() == '.' && isDigit(t.peekNext()) {
		t.advance()

		for isDigit(t.peek()) {
			t.advance()
		}
	}

	t.add(Token{
		Id: NUMERIC,
		Raw: RawToken{
			Value:  string([]rune(t.source)[t.beginOfLexeme:t.current]),
			Line:   t.line,
			Column: t.column,
		},
	})
}

func (t *Tokenizer) identifier() {
	for isAlphaNumeric(t.peek()) {
		t.advance()
	}

	keyword := getKeyword(string([]rune(t.source)[t.beginOfLexeme:t.current]))

	if keyword != UNKNOWN {
		t.addNonLiteralToken(keyword)
		return
	}

	t.add(Token{
		Id: IDENTIFIER,
		Raw: RawToken{
			Value:  string([]rune(t.source)[t.beginOfLexeme:t.current]),
			Line:   t.line,
			Column: t.column,
		},
	})
}

func (t *Tokenizer) scan() {
	char := t.advance()

	switch char {
	case '#':
		for !t.isAtEnd() && t.peek() != '\n' {
			t.advance()
		}
	case '(':
		t.addNonLiteralToken(LEFT_PARENTHESIS)
	case ')':
		t.addNonLiteralToken(RIGHT_PARENTHESIS)
	case '{':
		t.addNonLiteralToken(LEFT_BRACE)
	case '}':
		t.addNonLiteralToken(RIGHT_BRACE)
	case '[':
		t.addNonLiteralToken(LEFT_BRACKET)
	case ']':
		t.addNonLiteralToken(RIGHT_BRACKET)
	case ',':
		t.addNonLiteralToken(COMMA)
	case '.':
		t.addNonLiteralToken(DOT)
	case '-':
		t.addNonLiteralToken(MINUS)
	case '+':
		t.addNonLiteralToken(PLUS)
	case ':':
		t.addNonLiteralToken(COLON)
	case ';':
		t.addNonLiteralToken(SEMICOLON)
	case '*':
		t.addNonLiteralToken(STAR)
	case '/':
		t.addNonLiteralToken(SLASH)
	case '!':
		if t.match('=') {
			t.addNonLiteralToken(BANG_EQUAL)
		} else {
			t.addNonLiteralToken(NOT)
		}
	case '=':
		if t.match('=') {
			t.addNonLiteralToken(EQUAL_EQUAL)
		} else {
			t.addNonLiteralToken(EQUAL)
		}
	case '<':
		if t.match('=') {
			t.addNonLiteralToken(LESS_THAN_EQUAL)
		} else {
			t.addNonLiteralToken(LESS_THAN)
		}
	case '>':
		if t.match('=') {
			t.addNonLiteralToken(GREATER_THAN_EQUAL)
		} else {
			t.addNonLiteralToken(GREATER_THAN)
		}
	case ' ', '\r', '\t':
	case '\n':
		t.advanceLine()
	case '"':
		t.string()
	default:
		if isDigit(char) {
			t.numeric()
		} else if isAlpha(char) {
			t.identifier()
		} else {
			exception.NewSyntaxError("Unexpected character", t.line, t.column, string(char))
		}
	}
}

func Tokenize(source string) ([]Token, error) {
	tokenizer := Tokenizer{
		tokens:        []Token{},
		current:       0,
		beginOfLexeme: 0,
		line:          1,
		column:        0,
		source:        source,
	}

	for !tokenizer.isAtEnd() {
		tokenizer.beginOfLexeme = tokenizer.current

		tokenizer.scan()
	}

	tokenizer.add(Token{
		Id: EOF,
		Raw: RawToken{
			Value:  "",
			Line:   tokenizer.line,
			Column: tokenizer.column,
		},
	})

	return tokenizer.tokens, nil
}
