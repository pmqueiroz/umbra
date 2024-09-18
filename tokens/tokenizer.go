package tokens

import (
	"unicode"

	umbra_error "github.com/pmqueiroz/umbra/error"
)

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

	return append(tokens, Token{
		Id: EOF,
		Raw: RawToken{
			Value:  "",
			Line:   line,
			Column: column,
		},
	}), nil
}
