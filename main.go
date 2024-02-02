package main

import (
	"fmt"
)

func main() {
	dat, err := readFile("example.umb")

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	tokens, err := tokenizer(dat)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	for _, tok := range tokens {
		fmt.Printf("Token { type: '%s', value: '%s', start: %d, end: %d }\n", tok.id, tok.value, tok.line, tok.column)
	}
}
