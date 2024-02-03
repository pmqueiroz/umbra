package tokens

import (
	"fmt"
	"testing"
)

func TestGenerateTokenSingleCases(t *testing.T) {
	var tests = []struct {
		lexis string
		want  TokenType
	}{
		{"cleiton", IDENTIFIER},
		{".", PUNCTUATOR},
		{`"something string"`, STRING},
		{"true", BOOLEAN},
		{"false", BOOLEAN},
		{"123", NUMERIC},
		{"null", NULL},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("should return type %s when lexis is %s", testCase.lexis, testCase.want)
		t.Run(testName, func(t *testing.T) {
			result, err := generateToken(testCase.lexis, 0, 0)

			if err != nil {
				t.Fatal("return an unexpected error")
			}

			if result.Id != testCase.want {
				t.Errorf("got %s, want %s", result.Id, testCase.want)
			}
		})
	}
}
