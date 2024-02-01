package main

import (
	"fmt"
	"os"
)

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", &UmbraError{
			Code:    "MODULE_NOT_FOUND",
			Message: fmt.Sprintf("Cannot find module '%s'", path),
		}
	}

	return string(dat[:]), nil
}
