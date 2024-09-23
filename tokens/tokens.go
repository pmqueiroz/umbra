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
	LEFT_BRACKET       TokenType = "LEFT_BRACKET"
	RIGHT_BRACKET      TokenType = "RIGHT_BRACKET"
	COMMA              TokenType = "COMMA"
	DOT                TokenType = "DOT"
	MINUS              TokenType = "MINUS"
	PLUS               TokenType = "PLUS"
	COLON              TokenType = "COLON"
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
	STDOUT             TokenType = "STDOUT"
	STDERR             TokenType = "STDERR"
	RETURN             TokenType = "RETURN"
	THIS               TokenType = "THIS"
	TRUE               TokenType = "TRUE"
	FALSE              TokenType = "FALSE"
	CONST              TokenType = "VAR"
	MUT                TokenType = "MUT"
	STR_TYPE           TokenType = "STR_TYPE"
	NUM_TYPE           TokenType = "NUM_TYPE"
	BOOL_TYPE          TokenType = "BOOL_TYPE"
	VOID_TYPE          TokenType = "VOID_TYPE"
	ARR_TYPE           TokenType = "ARR_TYPE"
	HASHMAP_TYPE       TokenType = "HASHMAP_TYPE"
	MODULE             TokenType = "MODULE"
	BREAK              TokenType = "BREAK"
)

var DATA_TYPES = []TokenType{
	STR_TYPE,
	NUM_TYPE,
	BOOL_TYPE,
	HASHMAP_TYPE,
	ARR_TYPE,
}

var SPECIAL_TYPES = []TokenType{
	VOID_TYPE,
}

var reservedKeywordsMap = map[string]TokenType{
	"not":     NOT,
	"and":     AND,
	"else":    ELSE,
	"def":     FUN,
	"for":     FOR,
	"if":      IF,
	"null":    NULL,
	"or":      OR,
	"stdout":  STDOUT,
	"stderr":  STDERR,
	"return":  RETURN,
	"this":    THIS,
	"true":    TRUE,
	"false":   FALSE,
	"str":     STR_TYPE,
	"num":     NUM_TYPE,
	"bool":    BOOL_TYPE,
	"arr":     ARR_TYPE,
	"hashmap": HASHMAP_TYPE,
	"void":    VOID_TYPE,
	"const":   CONST,
	"mut":     MUT,
	"module":  MODULE,
	"break":   BREAK,
}

func getKeyword(lexis string) TokenType {
	if tokenType, exists := reservedKeywordsMap[lexis]; exists {
		return tokenType
	}

	return UNKNOWN
}

type ColumnRange struct {
	From int
	To   int
}

type Location struct {
	Line  int
	Range ColumnRange
}

type Token struct {
	Lexeme string
	Type   TokenType
	Loc    Location
}
