package exception

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

func NewGenericError(code string, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &GenericError{
		message: message,
		code:    code,
	}
}
