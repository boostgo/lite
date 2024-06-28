package to

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

func Bytes(anyValue any) []byte {
	return toBytes(anyValue, false)
}

func toBytes(anyValue any, memory bool) []byte {
	if anyValue == nil {
		return nil
	}

	valueType := reflect.TypeOf(anyValue)

	switch valueType.Kind() {
	case reflect.String:
		return StringToBytes(anyValue.(string))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return StringToBytes(toString(anyValue, false))
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		marshalled, err := json.Marshal(anyValue)
		if err != nil {
			return nil
		}

		return marshalled
	case reflect.Ptr:
		if memory {
			return nil
		}

		return toBytes(reflect.ValueOf(anyValue).Interface(), true)
	default:
		return StringToBytes(toString(anyValue, false))
	}
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}
