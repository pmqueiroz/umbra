package main

import (
	"fmt"
)

func main() {
	dat, err := readFile("example/example.umb")

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	fmt.Print(dat)
}
