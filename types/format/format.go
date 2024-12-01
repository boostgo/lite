package format

import (
	"strings"
	"unicode"
)

func Title(input string) string {
	if input == "" {
		return ""
	}

	input = strings.TrimSpace(input)
	runes := []rune(input)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}

	return string(runes)
}

func Username(input string) string {
	if input == "" {
		return ""
	}

	return strings.ToLower(
		strings.ReplaceAll(
			strings.TrimSpace(
				input,
			),
			" ", "_",
		),
	)
}

func Alpha(input string) string {
	if input == "" {
		return ""
	}

	var result []rune
	for _, r := range input {
		if !unicode.IsDigit(r) {
			result = append(result, r)
		}
	}

	return string(result)
}

func Numeric(input string) string {
	if input == "" {
		return ""
	}

	var result []rune
	for _, r := range input {
		if unicode.IsDigit(r) {
			result = append(result, r)
		}
	}

	return string(result)
}
