package ast

import (
	"fmt"

	"github.com/pmqueiroz/umbra/tokens"
)

func Parse(tokens []tokens.Token) {
	fmt.Println("Parsing...")
	for _, tok := range tokens {
		fmt.Printf(
			"Token { type: '%s', value: '%s', line: %d, column: %d }\n",
			tok.Id,
			tok.Raw.Value,
			tok.Raw.Line,
			tok.Raw.Column,
		)
	}
}
