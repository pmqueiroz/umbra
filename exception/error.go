package exception

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/pmqueiroz/umbra/globals"
)

type UmbraError struct {
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

func (e *UmbraError) header() string {
	errorType := ""
	switch {
	case strings.HasPrefix(e.code, "RT"):
		errorType = "RuntimeError"
	case strings.HasPrefix(e.code, "GN"):
		errorType = "GenericError"
	case strings.HasPrefix(e.code, "TY"):
		errorType = "TypeError"
	case strings.HasPrefix(e.code, "SY"):
		errorType = "SyntaxError"
	default:
		errorType = "UnknownError"
	}

	return color.New(color.Bold).Add(color.BgRed).Add(color.FgHiWhite).Sprintf("%s[%s]", errorType, e.code)
}

func annotation(e *UmbraError) string {
	annotation := e.node.Reference()
	lines := toLines(e.node.GetLocs())
	lines += " |"

	return fmt.Sprintf(`%s

%s %s
%s %s`, e.header(), color.BlueString(lines), annotation, strings.Repeat(" ", len(lines)+1)+color.RedString(strings.Repeat("^", len(annotation))), color.YellowString(e.message))
}

func (e *UmbraError) Error() string {
	if e.node != nil {
		return annotation(e)
	}

	return fmt.Sprintf("%s: %s", e.code, color.RedString(e.message))
}

func NewUmbraError(code string, node globals.Node, arguments ...any) error {
	message := fmt.Sprintf(Messages[code], arguments...)

	return &UmbraError{
		code:    code,
		message: message,
		node:    node,
	}
}
