package main

import (
	"fmt"

	"github.com/umbra-lang/umbra/tokens"
)

func main() {
	dat, err := readFile("example.umb")

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	tokens, err := tokens.Tokenizer(dat)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	for _, tok := range tokens {
		fmt.Printf("Token { type: '%s', value: '%s', line: %d, column: %d }\n", tok.Id, tok.Value, tok.Line, tok.Column)
	}
}
