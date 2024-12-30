package exception

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pmqueiroz/umbra/globals"
)

type RuntimeError struct {
	code    string
	message string
	node    globals.Node
}

func toLines(locs []globals.Loc) string {
	lines := []string{}

	for _, loc := range locs {
		lines = append(lines, strconv.Itoa(loc.Line))
	}

	slices.Sort(lines)
	return strings.Join(slices.Compact(lines), ",")
}

func (e *RuntimeError) Error() string {

	if e.node != nil {
		annotation := e.node.Reference()
		lines := toLines(e.node.GetLocs())
		lines += " |"

		return fmt.Sprintf(`%s

%s %s
%s %s`, color.New(color.Bold).Sprintf("RuntimeError[%s]", e.code), color.BlueString(lines), annotation, strings.Repeat(" ", len(lines)+1)+color.RedString(strings.Repeat("^", len(annotation))), color.YellowString(e.message))
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
