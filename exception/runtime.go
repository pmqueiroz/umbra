package exception

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pmqueiroz/umbra/globals"
)

type RuntimeError struct {
	code    string
	message string
	node    globals.Node
}

func (e *RuntimeError) Error() string {

	if e.node != nil {
		annotation := e.node.Reference()
		return fmt.Sprintf(`%s

%s
%s %s`, color.New(color.Bold).Sprintf("RuntimeError[%s]", e.code), annotation, color.RedString(strings.Repeat("^", len(annotation))), color.YellowString(e.message))
	}

	return fmt.Sprintf("RuntimeError[%s]: %s", e.code, color.RedString(e.message))
}

func NewRuntimeError(code string, node globals.Node, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &RuntimeError{
		code:    code,
		message: message,
		node:    node,
	}
}
