package to

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"unsafe"
)

// Bytes convert any value to bytes slice.
//
// If value is string calls BytesFromString function.
//
// If value is numeric convert it to string then to bytes.
//
// If value is uuid convert to string by String() function and then to bytes.
//
// If value is fmt.Stringer implementation calls .String() method and then to bytes.
func Bytes(value any) []byte {
	return toBytes(value, false)
}

func toBytes(value any, memory bool) []byte {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if v == nil {
			return nil
		}

		return v
	case string:
		return BytesFromString(v)
	case *string:
		if v == nil {
			return nil
		}

		return BytesFromString(*v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, bool:
		return BytesFromString(toString(v, false))
	case uuid.UUID:
		return BytesFromString(String(v))
	case *uuid.UUID:
		if v == nil {
			return nil
		}

		return BytesFromString(String(*v))
	case fmt.Stringer:
		if v == nil {
			return nil
		}

		return BytesFromString(v.String())
	}

	valueType := reflect.TypeOf(value)

	switch valueType.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		marshalled, err := json.Marshal(value)
		if err != nil {
			return nil
		}

		return marshalled
	case reflect.Ptr:
		if memory {
			return nil
		}

		return toBytes(reflect.ValueOf(value).Interface(), true)
	default:
		return BytesFromString(toString(value, false))
	}
}

// BytesFromString converts string to bytes slice with no allocation
func BytesFromString(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}
