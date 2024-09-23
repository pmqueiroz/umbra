package tokens

import (
	"fmt"
	"testing"
)

func TestGenerateTokenSingleCases(t *testing.T) {
	var tests = []struct {
		source string
		want   TokenType
	}{
		{"cleiton", IDENTIFIER},
		{".", DOT},
		{`"something string"`, STRING},
		{"true", TRUE},
		{"false", FALSE},
		{"123", NUMERIC},
		{"null", NULL},
	}

	for _, testCase := range tests {
		testName := fmt.Sprintf("should return type %s when lexis is %s", testCase.source, testCase.want)
		t.Run(testName, func(t *testing.T) {
			result, err := Tokenize(testCase.source)

			if err != nil || len(result) != 1 {
				t.Fatal("return an unexpected error")
			}

			if result[0].Type != testCase.want {
				t.Errorf("got %s, want %s", result[0].Type, testCase.want)
			}
		})
	}
}
