package tokens

import "unicode"

func isKeyword(lexis string) bool {
	for _, keyword := range reservedKeywords {
		if keyword == Keyword(lexis) {
			return true
		}
	}

	return false
}

func isIsolatedPunctuator(lexis string) bool {
	for _, punctuator := range isolatedPunctuators {
		if punctuator == Punctuator(lexis) {
			return true
		}
	}

	return false
}

func isPunctuator(lexis string) bool {
	if isIsolatedPunctuator(lexis) {
		return true
	}

	for _, punctuator := range combinedPunctuators {
		if punctuator == Punctuator(lexis) {
			return true
		}
	}

	return false
}

func isValidString(lexis string) bool {
	if len(lexis) >= 2 && lexis[len(lexis)-1] == '"' {
		return true
	}
	return false
}

func isValidNumeric(lexis string) bool {
	// TODO: add checks for floats
	for _, runes := range lexis {
		if !unicode.IsDigit(runes) {
			return false
		}
	}

	return true
}

func isBoolean(lexis string) bool {
	return lexis == "true" || lexis == "false"
}
