package error

import (
	"fmt"
)

type GenericError struct {
	code    string
	message string
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("Error: %s %s", e.message, e.prettyCtx())
}

func (e *GenericError) prettyCtx() string {
	return fmt.Sprintf(`{
   code: '%s'		
}`, e.code)
}

func NewGenericError(code string, message string) error {
	return &GenericError{
		message: message,
		code:    code,
	}
}

type SyntaxError struct {
	message string
	line    int
	column  int
	raw     string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("SyntaxError: %s %s", e.message, e.prettyCtx())
}

func (e *SyntaxError) prettyCtx() string {
	return fmt.Sprintf(`{
   line: %d		
   column: %d		
   raw: %s		
}`, e.line, e.column, e.raw)
}

func NewSyntaxError(message string, line int, column int, raw string) error {
	return &SyntaxError{
		message: message,
		line:    line,
		column:  column,
		raw:     raw,
	}
}
