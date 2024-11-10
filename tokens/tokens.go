package tokens

type TokenType string

const (
	UNKNOWN            TokenType = "UNKNOWN"
	EOF                TokenType = "EOF"
	IDENTIFIER         TokenType = "IDENTIFIER"
	STRING             TokenType = "STRING"
	CHAR               TokenType = "CHAR"
	NUMERIC            TokenType = "NUMERIC"
	LEFT_PARENTHESIS   TokenType = "LEFT_PAREN"
	RIGHT_PARENTHESIS  TokenType = "RIGHT_PAREN"
	LEFT_BRACE         TokenType = "LEFT_BRACE"
	RIGHT_BRACE        TokenType = "RIGHT_BRACE"
	LEFT_BRACKET       TokenType = "LEFT_BRACKET"
	RIGHT_BRACKET      TokenType = "RIGHT_BRACKET"
	COMMA              TokenType = "COMMA"
	DOT                TokenType = "DOT"
	VARIADIC           TokenType = "VARIADIC"
	MINUS              TokenType = "MINUS"
	PLUS               TokenType = "PLUS"
	COLON              TokenType = "COLON"
	DOUBLE_COLON       TokenType = "DOUBLE_COLON"
	SEMICOLON          TokenType = "SEMICOLON"
	SLASH              TokenType = "SLASH"
	PERCENT            TokenType = "PERCENT"
	STAR               TokenType = "STAR"
	EQUAL              TokenType = "EQUAL"
	EQUAL_EQUAL        TokenType = "EQUAL_EQUAL"
	BANG_EQUAL         TokenType = "BANG_EQUAL"
	GREATER_THAN       TokenType = "GREATER_THAN"
	GREATER_THAN_EQUAL TokenType = "GREATER_THAN_EQUAL"
	LESS_THAN          TokenType = "LESS_THAN"
	LESS_THAN_EQUAL    TokenType = "LESS_THAN_EQUAL"
	TILDE              TokenType = "TILDE"
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
	CHAR_TYPE          TokenType = "CHAR_TYPE"
	NUM_TYPE           TokenType = "NUM_TYPE"
	BOOL_TYPE          TokenType = "BOOL_TYPE"
	VOID_TYPE          TokenType = "VOID_TYPE"
	ARR_TYPE           TokenType = "ARR_TYPE"
	HASHMAP_TYPE       TokenType = "HASHMAP_TYPE"
	FUN_TYPE           TokenType = "FUN_TYPE"
	ANY_TYPE           TokenType = "ANY_TYPE"
	BREAK              TokenType = "BREAK"
	CONTINUE           TokenType = "CONTINUE"
	PUBLIC             TokenType = "PUBLIC"
	IMPORT             TokenType = "IMPORT"
	HOOK               TokenType = "HOOK"
	NOT_A_NUMBER       TokenType = "NOT_A_NUMBER"
	RANGE              TokenType = "RANGE"
)

var DATA_TYPES = []TokenType{
	STR_TYPE,
	CHAR_TYPE,
	NUM_TYPE,
	BOOL_TYPE,
	HASHMAP_TYPE,
	ARR_TYPE,
	ANY_TYPE,
	FUN_TYPE,
}

var SPECIAL_TYPES = []TokenType{
	VOID_TYPE,
}

var reservedKeywordsMap = map[string]TokenType{
	"not":      NOT,
	"and":      AND,
	"else":     ELSE,
	"def":      FUN,
	"for":      FOR,
	"if":       IF,
	"null":     NULL,
	"or":       OR,
	"stdout":   STDOUT,
	"stderr":   STDERR,
	"return":   RETURN,
	"this":     THIS,
	"true":     TRUE,
	"false":    FALSE,
	"str":      STR_TYPE,
	"char":     CHAR_TYPE,
	"num":      NUM_TYPE,
	"bool":     BOOL_TYPE,
	"arr":      ARR_TYPE,
	"hashmap":  HASHMAP_TYPE,
	"void":     VOID_TYPE,
	"func":     FUN_TYPE,
	"const":    CONST,
	"mut":      MUT,
	"break":    BREAK,
	"continue": CONTINUE,
	"pub":      PUBLIC,
	"import":   IMPORT,
	"any":      ANY_TYPE,
	"NaN":      NOT_A_NUMBER,
	"range":    RANGE,
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
