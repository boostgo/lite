package to

import "strings"

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
