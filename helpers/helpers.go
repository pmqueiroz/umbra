package helpers

import (
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func ReadFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewUmbraError("GN001", nil, path)
	}

	return string(dat[:]), nil
}
