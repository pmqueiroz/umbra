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
	"if",
	"else",
}

var isolatedPunctuators = [...]Punctuator{
	"=",
	"{",
	"}",
	"(",
	")",
	"[",
	"]",
	".",
	":",
	",",
	"+",
	"-",
	"/",
	"*",
}

var combinedPunctuators = [...]Punctuator{
	":=",
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

func isIsolatedPunctuator(lexis string) bool {
	for _, punctuator := range isolatedPunctuators {
		if punctuator == Punctuator(lexis) {
			return true
		}
	}

	return false
}

func isPunctuator(lexis string) bool {
	if isIsolatedPunctuator(lexis) {
		return true
	}

	for _, punctuator := range combinedPunctuators {
		if punctuator == Punctuator(lexis) {
			return true
		}
	}

	return false
}

func isValidString(lexis string) bool {
	if len(lexis) >= 2 && lexis[len(lexis)-1] == '"' {
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

	for index, char := range code {
		if comment {
			if char == '\n' {
				line++
				column = 1
				comment = false
			}

			continue
		}

		if char == '"' {
			var closingQuote = stringInProgress && !(code[index-1] == '\\') // is closing quote if isn´t escaped

			if !stringInProgress {
				if lexis != "" {
					token, err := generateToken(lexis, line, column-len(lexis))
					if err != nil {
						return tokens, err
					}
					tokens = append(tokens, token)
					lexis = ""
				}

				stringInProgress = true
			}

			if closingQuote {
				stringInProgress = false

				token, err := generateToken(lexis+string(char), line, column-len(lexis))
				if err != nil {
					return tokens, err
				}
				tokens = append(tokens, token)
				lexis = ""
				column++
				continue
			}
		}

		if char == '#' && !stringInProgress {
			if !comment {
				comment = true
			}
			continue
		}

		if (unicode.IsPunct(char) || unicode.IsSymbol(char)) && !stringInProgress {
			if isIsolatedPunctuator(string(char)) {
				if lexis != "" {
					token, err := generateToken(lexis, line, column-len(lexis))
					if err != nil {
						return tokens, err
					}

					isolatedToken, err := generateToken(string(char), line, column)
					if err != nil {
						return tokens, err
					}
					tokens = append(tokens, token, isolatedToken)
					lexis = ""
					column++
					continue
				} else {
					token, err := generateToken(string(char), line, column-len(lexis))
					if err != nil {
						return tokens, err
					}
					tokens = append(tokens, token)
					lexis = ""
					column++
					continue
				}
			}
		}

		if unicode.IsSpace(char) {
			if lexis != "" && !stringInProgress {
				token, err := generateToken(lexis, line, column-len(lexis))
				if err != nil {
					return tokens, err
				}
				tokens = append(tokens, token)
				lexis = ""
			}

			column++

			if char == '\n' {
				if stringInProgress {
					return tokens, &SyntaxError{
						message: "Breaking the line before terminating the string",
						line:    line,
						column:  column,
						raw:     lexis,
					}
				}

				line++
				column = 1
			}

			if stringInProgress {
				lexis += string(char)
				column++
			}

			continue
		}

		lexis += string(char)
		column++
	}

	if len(lexis) != 0 {
		return tokens, &UmbraError{
			Code:    "INTERNAL_ERROR",
			Message: "Could not resolve entire file",
		}
	}

	return tokens, nil
}
