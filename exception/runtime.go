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

func NewRuntimeError(message string) error {
	return &RuntimeError{
		message: message,
	}
}
