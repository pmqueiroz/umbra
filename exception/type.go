package exception

import (
	"fmt"
)

type TypeError struct {
	message string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("TypeError: %s", e.message)
}

func NewTypeError(message string) error {
	return &TypeError{
		message: message,
	}
}
