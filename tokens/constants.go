package tokens

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
	EOF        TokenType = "EOF"
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
