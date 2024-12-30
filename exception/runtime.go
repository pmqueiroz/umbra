package exception

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type RuntimeError struct {
	code       string
	message    string
	annotation string
}

func (e *RuntimeError) Error() string {

	if e.annotation != "" {
		return fmt.Sprintf(`%s

%s
%s %s`, color.New(color.Bold).Sprintf("RuntimeError[%s]", e.code), e.annotation, color.RedString(strings.Repeat("^", len(e.annotation))), color.YellowString(e.message))
	}

	return fmt.Sprintf("RuntimeError[%s]: %s", e.code, e.message)
}

func NewRuntimeError(code string, annotation string, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &RuntimeError{
		code:       code,
		message:    message,
		annotation: annotation,
	}
}
