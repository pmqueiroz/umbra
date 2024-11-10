package exception

import (
	"fmt"
)

type RuntimeError struct {
	code    string
	message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError[%s]: %s", e.code, e.message)
}

func NewRuntimeError(code string, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &RuntimeError{
		code:    code,
		message: message,
	}
}
