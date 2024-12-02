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

// Title format input to "title" format.
// "Title" format is first sentence letter uppercase and other lowercase.
// Example:
//
//	Input: Hello WORLD
//	Output: Hello world
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

// Code format input to "code" format.
// "Code" format is lowercase and no-space text.
// Example:
//
//	Input: Hello-World. 123 !!! 777
//	Output: hello_world_123_777
func Code(input string) string {
	if input == "" {
		return ""
	}

	return strings.ReplaceAll(strings.ToLower(AlphaNumeric(input)), " ", "_")
}

// Alpha format input to only latin-letter text
// Example:
//
//	Input: Hello-World. 123 !!!
//	Output: Hello World
func Alpha(input string) string {
	if input == "" {
		return ""
	}

	var result []rune
	for _, r := range input {
		if isLatin(r) || r == ' ' {
			result = append(result, r)
		} else {
			result = append(result, ' ')
		}
	}

	return clearInput(string(result))
}

// Numeric format input to only digit text
// Example:
//
//	Input: Hello-World. 123 !!!
//	Output: 123
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

	return clearInput(string(result))
}

// AlphaNumeric format input to text with latin letters & digits (other symbols - erase).
// Example:
//
//	Input: Hello-World. 123 !!!
//	Output: Hello World 123
func AlphaNumeric(input string) string {
	if input == "" {
		return ""
	}

	var result []rune
	for _, r := range input {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || r == ' ' {
			result = append(result, r)
		} else {
			result = append(result, ' ')
		}
	}

	return clearInput(string(result))
}

// Cyrillic format input from cyrillic text to latin-code format.
// Example:
//
//	Input: Привет
//	Output: privet
func Cyrillic(input string) string {
	if input == "" {
		return ""
	}

	return Code(iuliia.Wikipedia.Translate(clearInput(input)))
}

// EveryTitle makes every word start with uppercase.
// Example:
//
//	Input: HELLO WORLD
//	Output: Hello World
func EveryTitle(input string) string {
	return cases.Title(language.Und).String(strings.ToLower(clearInput(input)))
}

// Name format input to First/Last name format.
// Example:
//
//	Input: john smith
//	Output: John Smith
func Name(input string) string {
	return EveryTitle(Alpha(input))
}

func clearInput(input string) string {
	return _manySpacesReg.ReplaceAllString(strings.TrimSpace(input), " ")
}

func isLatin(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
