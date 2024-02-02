package main

import (
	"fmt"
	"os"

	umbra_error "github.com/umbra-lang/umbra/error"
)

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", umbra_error.NewGenericError("MODULE_NOT_FOUND", fmt.Sprintf("Cannot find module '%s'", path))
	}

	return string(dat[:]), nil
}
