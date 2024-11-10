package cli

import (
	"github.com/pmqueiroz/umbra/ast"
	"github.com/pmqueiroz/umbra/tokens"
	"github.com/sanity-io/litter"
)

func PrintAst(module ast.ModuleStatement) {
	litter.Dump(module)
}

func PrintTokens(tokens []tokens.Token) {
	litter.Dump(tokens)
}
