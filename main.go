package main

import (
	"fmt"
)

func main() {
	dat, err := readFile("example/example.umb")

	if err != nil {
		panic(err)
	}

	fmt.Print(dat)
}
