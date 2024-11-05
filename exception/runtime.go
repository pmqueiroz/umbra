package exception

import (
	"fmt"
)

type RuntimeError struct {
	message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s", e.message)
}

func NewRuntimeError(code string, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &RuntimeError{
		message: message,
	}
}
