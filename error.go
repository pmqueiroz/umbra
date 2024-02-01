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
