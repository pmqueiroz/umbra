package main

import (
	"os"
)

func readFile(path string) (string, error) {
	dat, err := os.ReadFile("example/example.umb")

	if err != nil {
		return "", err
	}

	return string(dat[:]), nil
}
