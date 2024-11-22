package to

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"unsafe"
)

func Bytes(value any) []byte {
	return toBytes(value, false)
}

func toBytes(value any, memory bool) []byte {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		return BytesFromString(v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, bool:
		return BytesFromString(toString(v, false))
	case uuid.UUID:
		return BytesFromString(v.String())
	case fmt.Stringer:
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

func BytesFromString(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}
