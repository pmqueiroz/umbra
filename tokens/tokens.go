package tokens

type TokenType string

const (
	UNKNOWN            TokenType = "UNKNOWN"
	EOF                TokenType = "EOF"
	IDENTIFIER         TokenType = "IDENTIFIER"
	STRING             TokenType = "STRING"
	NUMERIC            TokenType = "NUMERIC"
	LEFT_PARENTHESIS   TokenType = "LEFT_PAREN"
	RIGHT_PARENTHESIS  TokenType = "RIGHT_PAREN"
	LEFT_BRACE         TokenType = "LEFT_BRACE"
	RIGHT_BRACE        TokenType = "RIGHT_BRACE"
	COMMA              TokenType = "COMMA"
	DOT                TokenType = "DOT"
	MINUS              TokenType = "MINUS"
	PLUS               TokenType = "PLUS"
	SEMICOLON          TokenType = "SEMICOLON"
	SLASH              TokenType = "SLASH"
	STAR               TokenType = "STAR"
	EQUAL              TokenType = "EQUAL"
	EQUAL_EQUAL        TokenType = "EQUAL_EQUAL"
	BANG_EQUAL         TokenType = "BANG_EQUAL"
	GREATER_THAN       TokenType = "GREATER_THAN"
	GREATER_THAN_EQUAL TokenType = "GREATER_THAN_EQUAL"
	LESS_THAN          TokenType = "LESS_THAN"
	LESS_THAN_EQUAL    TokenType = "LESS_THAN_EQUAL"
	NOT                TokenType = "NOT"
	AND                TokenType = "AND"
	ELSE               TokenType = "ELSE"
	FUN                TokenType = "FUN"
	FOR                TokenType = "FOR"
	IF                 TokenType = "IF"
	NULL               TokenType = "NULL"
	OR                 TokenType = "OR"
	PRINT              TokenType = "PRINT"
	RETURN             TokenType = "RETURN"
	THIS               TokenType = "THIS"
	TRUE               TokenType = "TRUE"
	FALSE              TokenType = "FALSE"
	STR_VAR            TokenType = "STR_VAR"
	ARR_VAR            TokenType = "ARR_VAR"
	OBJ_VAR            TokenType = "OBJ_VAR"
	NUM_VAR            TokenType = "NUM_VAR"
	WHILE              TokenType = "WHILE"
)

var reservedKeywordsMap = map[string]TokenType{
	"not":    NOT,
	"and":    AND,
	"else":   ELSE,
	"def":    FUN,
	"for":    FOR,
	"if":     IF,
	"null":   NULL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"this":   THIS,
	"true":   TRUE,
	"false":  FALSE,
	"str":    STR_VAR,
	"arr":    ARR_VAR,
	"obj":    OBJ_VAR,
	"num":    NUM_VAR,
	"while":  WHILE,
}

func getKeyword(lexis string) TokenType {
	if tokenType, exists := reservedKeywordsMap[lexis]; exists {
		return tokenType
	}

	return UNKNOWN
}

type RawToken struct {
	Value  string
	Line   int
	Column int
}

type Token struct {
	Raw RawToken
	Id  TokenType
}
