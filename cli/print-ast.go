package cli

import (
	"github.com/pmqueiroz/umbra/ast"
	"github.com/sanity-io/litter"
)

func PrintAst(module ast.ModuleStatement) {
	litter.Dump(module)
}
