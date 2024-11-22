package to

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
	case uuid.UUID:
		var buf [36]byte
		encodeUuidHex(buf[:], v)
		return StringFromBytes(buf[:])
	case fmt.Stringer:
		return v.String()
	case []byte:
		return StringFromBytes(v)
	case error:
		if v == nil {
			return ""
		}

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
	case rune:
		return string(v)
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
		return StringFromBytes(valueInBytes)
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

func StringFromBytes(buffer []byte) string {
	return *(*string)(unsafe.Pointer(&buffer))
}

func encodeUuidHex(dst []byte, uuid uuid.UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}
