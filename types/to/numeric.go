package to

import (
	"strconv"
)

// Int convert any value to int.
//
// If value is nil return 0.
func Int(anyValue any) int {
	if anyValue == nil {
		return 0
	}

	if ptrValue, isPtr := anyValue.(*int); isPtr {
		return *ptrValue
	}

	switch value := anyValue.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case uint:
		return int(value)
	case uint8:
		return int(value)
	case uint16:
		return int(value)
	case uint32:
		return int(value)
	case uint64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case *int:
		if value == nil {
			return 0
		}

		return *value
	case *int8:
		if value == nil {
			return 0
		}

		return int(*value)
	case *int16:
		if value == nil {
			return 0
		}

		return int(*value)
	case *int32:
		if value == nil {
			return 0
		}

		return int(*value)
	case *int64:
		if value == nil {
			return 0
		}

		return int(*value)
	case *uint:
		if value == nil {
			return 0
		}

		return int(*value)
	case *uint8:
		if value == nil {
			return 0
		}

		return int(*value)
	case *uint16:
		if value == nil {
			return 0
		}

		return int(*value)
	case *uint32:
		if value == nil {
			return 0
		}

		return int(*value)
	case *uint64:
		if value == nil {
			return 0
		}

		return int(*value)
	case *float32:
		if value == nil {
			return 0
		}

		return int(*value)
	case *float64:
		if value == nil {
			return 0
		}

		return int(*value)
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}

		return parsed
	default:
		stringValue := String(value)
		if stringValue == "" {
			return 0
		}

		parsed, err := strconv.Atoi(stringValue)
		if err != nil {
			return 0
		}

		return parsed
	}
}

// Float32 convert any value to float32.
//
// If value is nil return 0.
func Float32(anyValue any) float32 {
	if anyValue == nil {
		return 0
	}

	switch value := anyValue.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case *float32:
		if value == nil {
			return 0
		}

		return *value
	case int:
		return float32(value)
	case int8:
		return float32(value)
	case int16:
		return float32(value)
	case int32:
		return float32(value)
	case int64:
		return float32(value)
	case uint:
		return float32(value)
	case uint8:
		return float32(value)
	case uint16:
		return float32(value)
	case uint32:
		return float32(value)
	case uint64:
		return float32(value)
	case string:
		parsed, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return 0
		}

		return float32(parsed)
	default:
		strValue := String(value)
		if strValue == "" {
			return 0
		}

		parsed, err := strconv.ParseFloat(strValue, 32)
		if err != nil {
			return 0
		}

		return float32(parsed)
	}
}

// Float64 convert any value to float64.
//
// If value is nil return 0.
func Float64(anyValue any) float64 {
	if anyValue == nil {
		return 0
	}

	switch value := anyValue.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case *float64:
		if value == nil {
			return 0
		}

		return *value
	case int:
		return float64(value)
	case int8:
		return float64(value)
	case int16:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case uint:
		return float64(value)
	case uint8:
		return float64(value)
	case uint16:
		return float64(value)
	case uint32:
		return float64(value)
	case uint64:
		return float64(value)
	case string:
		parsed, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}

		return parsed
	default:
		strValue := String(value)
		if strValue == "" {
			return 0
		}

		parsed, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return 0
		}

		return parsed
	}
}
