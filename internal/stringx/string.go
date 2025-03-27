package stringx

import "strings"

// Replace by provided map
func Replace(input string, replacers map[string]string) string {
	for oldValue, newValue := range replacers {
		input = strings.ReplaceAll(input, oldValue, newValue)
	}

	return input
}
