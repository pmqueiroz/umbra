package exception

import (
	"fmt"
)

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
