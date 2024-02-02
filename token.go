package main

import (
	"fmt"
	"unicode"
)

type TokenType string
type Keyword string
type Punctuator string

const (
	UNKNOWN    TokenType = "UNKNOWN"
	KEYWORD    TokenType = "KEYWORD"
	IDENTIFIER TokenType = "IDENTIFIER"
	PUNCTUATOR TokenType = "PUNCTUATOR"
	STRING     TokenType = "STRING"
	BOOLEAN    TokenType = "BOOLEAN"
	NUMERIC    TokenType = "NUMERIC"
	NULL       TokenType = "NULL"
)

var reservedKeywords = [...]Keyword{
	"package",
	"str",
	"mut",
	"obj",
	"arr",
	"def",
}

var validPunctuators = [...]Punctuator{
	"=",
	":=",
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
}

type Token struct {
	id     TokenType
	value  string
	line   int
	column int
}

func isKeyword(lexis string) bool {
	for _, keyword := range reservedKeywords {
		if keyword == Keyword(lexis) {
			return true
		}
	}

	return false
}

func isPunctuator(lexis string) bool {
	for _, punctuator := range validPunctuators {
		if punctuator == Punctuator(lexis) {
			return true
		}
	}

	return false
}

func isValidString(lexis string) bool {
	if lexis[len(lexis)-1] == '"' {
		return true
	}
	return false
}

func isValidNumeric(lexis string) bool {
	// TODO: add checks for floats
	for _, runes := range lexis {
		if !unicode.IsDigit(runes) {
			return false
		}
	}

	return true
}

func generateToken(lexis string, line int, column int) (Token, error) {
	if unicode.IsLetter([]rune(lexis)[0]) {
		if isKeyword(lexis) {
			return Token{
				id:     KEYWORD,
				value:  lexis,
				line:   line,
				column: column,
			}, nil
		}

		return Token{
			id:     IDENTIFIER,
			value:  lexis,
			line:   line,
			column: column,
		}, nil
	}

	if isPunctuator(lexis) {
		return Token{
			id:     PUNCTUATOR,
			value:  lexis,
			line:   line,
			column: column,
		}, nil
	}

	if lexis[0] == '"' {
		if isValidString(lexis) {
			return Token{
				id:     STRING,
				value:  lexis,
				line:   line,
				column: column,
			}, nil
		} else {
			return Token{
					id:     UNKNOWN,
					value:  lexis,
					line:   line,
					column: column,
				}, &SyntaxError{
					message: "Unterminated string",
					line:    line,
					column:  column,
					raw:     lexis,
				}
		}
	}

	if unicode.IsNumber([]rune(lexis)[0]) {
		if isValidNumeric(lexis) {
			return Token{
				id:     NUMERIC,
				value:  lexis,
				line:   line,
				column: column,
			}, nil
		} else {
			return Token{
					id:     UNKNOWN,
					value:  lexis,
					line:   line,
					column: column,
				}, &SyntaxError{
					message: "Invalid number",
					line:    line,
					column:  column,
				}
		}
	}

	return Token{
			id:     UNKNOWN,
			value:  lexis,
			line:   line,
			column: column,
		}, &SyntaxError{
			message: fmt.Sprintf("Unexpected lexis '%s'", lexis),
			line:    line,
			column:  column,
			raw:     lexis,
		}
}

func tokenizer(code string) ([]Token, error) {
	var tokens []Token
	var lexis string
	var line, column int = 1, 1
	var comment bool
	var stringInProgress bool

	for _, char := range code {
		if char == '#' {
			comment = true
			continue
		}

		if unicode.IsSpace(char) || comment || stringInProgress {
			if lexis != "" {
				token, err := generateToken(lexis, line, column-len(lexis))

				if err != nil {
					return tokens, err
				}

				tokens = append(tokens, token)
				lexis = ""
			}

			column++

			if char == '\n' {
				line++
				column = 1
				comment = false
			}

			continue
		}

		// if char == '"' {
		// 	var closingQuote = !stringInProgress && !(code[index-1] == '\\') // is closing quote if isn´t escaped

		// 	if !stringInProgress || closingQuote {

		// 		// terminar anterior

		// 		if !stringInProgress {
		// 			stringInProgress = true
		// 		} else if closingQuote {
		// 			stringInProgress = false
		// 		}

		// 	}
		// }

		lexis += string(char)
		column++
	}

	return tokens, nil
}
