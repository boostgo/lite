package to

import "strings"

// Bool convert any value to bool.
// If value is string, convert it to string and then compare for "true" value.
// If value is numeric and values equals to 1 then it's true.
// Other cases convert to string and compare to "true"
func Bool(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return v == 1
	case string:
		return strings.ToLower(v) == "true"
	default:
		return strings.ToLower(String(value)) == "true"
	}
}
