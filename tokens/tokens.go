package tokens

import (
	"fmt"
	"unicode"

	umbra_error "github.com/pmqueiroz/umbra/error"
)

type RawToken struct {
	Value  string
	Line   int
	Column int
}

type Token struct {
	Raw RawToken
	Id  TokenType
}

func generateToken(lexis string, line int, column int) (Token, error) {
	rawToken := RawToken{
		Value:  lexis,
		Line:   line,
		Column: column,
	}

	if unicode.IsLetter([]rune(lexis)[0]) {
		if isKeyword(lexis) {
			return Token{
				Id:  KEYWORD,
				Raw: rawToken,
			}, nil
		}

		if isBoolean(lexis) {
			return Token{
				Id:  BOOLEAN,
				Raw: rawToken,
			}, nil
		}

		if lexis == "null" {
			return Token{
				Id:  NULL,
				Raw: rawToken,
			}, nil
		}

		return Token{
			Id:  IDENTIFIER,
			Raw: rawToken,
		}, nil
	}

	if isPunctuator(lexis) {
		return Token{
			Id:  PUNCTUATOR,
			Raw: rawToken,
		}, nil
	}

	if lexis[0] == '"' {
		if isValidString(lexis) {
			return Token{
				Id:  STRING,
				Raw: rawToken,
			}, nil
		} else {
			return Token{
					Id:  UNKNOWN,
					Raw: rawToken,
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
				Id:  NUMERIC,
				Raw: rawToken,
			}, nil
		} else {
			return Token{
				Id:  UNKNOWN,
				Raw: rawToken,
			}, umbra_error.NewSyntaxError("Invalid number", line, column, lexis)
		}
	}

	return Token{
		Id:  UNKNOWN,
		Raw: rawToken,
	}, umbra_error.NewSyntaxError(fmt.Sprintf("Unexpected lexis '%s'", lexis), line, column, lexis)
}
