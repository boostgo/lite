package to

import (
	"encoding/json"
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
		return StringToBytes(v)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64, bool:
		return StringToBytes(toString(v, false))
	case uuid.UUID:
		return StringToBytes(v.String())
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
		return StringToBytes(toString(value, false))
	}
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}
