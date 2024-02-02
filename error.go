package main

import (
	"fmt"
)

type UmbraError struct {
	Code    string
	Message string
}

func (e *UmbraError) Error() string {
	return fmt.Sprintf("Error: %s %s", e.Message, e.prettyCtx())
}

func (e *UmbraError) prettyCtx() string {
	return fmt.Sprintf(`{
   code: '%s'		
}`, e.Code)
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
