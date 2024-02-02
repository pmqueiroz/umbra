package tokens

import (
	"fmt"
	"unicode"

	umbra_error "github.com/umbra-lang/umbra/error"
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
	Id     TokenType
	Value  string
	Line   int
	Column int
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
				Id:     KEYWORD,
				Value:  lexis,
				Line:   line,
				Column: column,
			}, nil
		}

		return Token{
			Id:     IDENTIFIER,
			Value:  lexis,
			Line:   line,
			Column: column,
		}, nil
	}

	if isPunctuator(lexis) {
		return Token{
			Id:     PUNCTUATOR,
			Value:  lexis,
			Line:   line,
			Column: column,
		}, nil
	}

	if lexis[0] == '"' {
		if isValidString(lexis) {
			return Token{
				Id:     STRING,
				Value:  lexis,
				Line:   line,
				Column: column,
			}, nil
		} else {
			return Token{
					Id:     UNKNOWN,
					Value:  lexis,
					Line:   line,
					Column: column,
				}, umbra_error.NewSyntaxError(
					"Unterminated string",
					line,
					column,
					lexis,
				)
		}
	}

	if unicode.IsNumber([]rune(lexis)[0]) {
		if isValidNumeric(lexis) {
			return Token{
				Id:     NUMERIC,
				Value:  lexis,
				Line:   line,
				Column: column,
			}, nil
		} else {
			return Token{
				Id:     UNKNOWN,
				Value:  lexis,
				Line:   line,
				Column: column,
			}, umbra_error.NewSyntaxError("Invalid number", line, column, lexis)
		}
	}

	return Token{
		Id:     UNKNOWN,
		Value:  lexis,
		Line:   line,
		Column: column,
	}, umbra_error.NewSyntaxError(fmt.Sprintf("Unexpected lexis '%s'", lexis), line, column, lexis)
}

func Tokenizer(code string) ([]Token, error) {
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
					return tokens, umbra_error.NewSyntaxError("Breaking the line before terminating the string", line, column, lexis)
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
		return tokens, umbra_error.NewGenericError("INTERNAL_ERROR", "Could not resolve entire file")
	}

	return tokens, nil
}
