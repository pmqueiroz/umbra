package main

import (
	"fmt"
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func readFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewGenericError("MODULE_NOT_FOUND", fmt.Sprintf("Cannot find module '%s'", path))
	}

	return string(dat[:]), nil
}
