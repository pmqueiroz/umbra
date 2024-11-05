package helpers

import (
	"os"

	"github.com/pmqueiroz/umbra/exception"
)

func ReadFile(path string) (string, error) {
	dat, err := os.ReadFile(path)

	if err != nil {
		return "", exception.NewGenericError("GN001", path)
	}

	return string(dat[:]), nil
}
