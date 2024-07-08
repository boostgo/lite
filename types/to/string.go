package to

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

func String(value any) string {
	return toString(value, false)
}

func toString(value any, memory bool) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case fmt.Stringer:
		return v.String()
	case []byte:
		return BytesToString(v)
	case error:
		return v.Error()
	case *string:
		if v == nil {
			return ""
		}

		return *v
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	}

	valueReflect := reflect.ValueOf(value)

	switch valueReflect.Kind() {
	case reflect.Ptr:
		if memory || valueReflect.IsZero() {
			return fmt.Sprintf("%v", value)
		}
		return toString(valueReflect.Elem().Interface(), true)
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		valueInBytes, err := json.Marshal(value)
		if err != nil {
			return ""
		}
		return string(valueInBytes)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(valueReflect.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(valueReflect.Uint(), 10)
	case reflect.Float32:
		return fmt.Sprintf("%f", valueReflect.Float())
	case reflect.Float64:
		return fmt.Sprintf("%g", valueReflect.Float())
	default:
		return fmt.Sprintf("%v", value)
	}
}

func BytesToString(buffer []byte) string {
	return *(*string)(unsafe.Pointer(&buffer))
}
