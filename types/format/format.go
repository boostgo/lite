package format

import (
	"github.com/mehanizm/iuliia-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
	"unicode"
)

var (
	_manySpacesReg = regexp.MustCompile("\\s{2,}")
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

	return clearInput(string(runes))
}

func Code(input string) string {
	if input == "" {
		return ""
	}

	return strings.ReplaceAll(strings.ToLower(clearInput(input)), " ", "_")
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

	return clearInput(string(result))
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

func Cyrillic(input string) string {
	if input == "" {
		return ""
	}

	return Code(iuliia.Wikipedia.Translate(clearInput(input)))
}

// EveryTitle makes every word start with uppercase
// Example:
//
//	Input: HELLO WORLD
//	Output: Hello World
func EveryTitle(text string) string {
	return cases.Title(language.Und).String(strings.ToLower(clearInput(text)))
}

func Name(name string) string {
	name = clearInput(name)
	name = EveryTitle(name)
	name = Alpha(name)
	return name
}

func clearInput(input string) string {
	return _manySpacesReg.ReplaceAllString(strings.TrimSpace(input), " ")
}
