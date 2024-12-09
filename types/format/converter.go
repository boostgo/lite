package format

import (
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/try"
	"github.com/boostgo/lite/types/flex"
	"math"
	"reflect"
	"strings"
)

type StringFormatter func(string) string
type IntFormatter func(int) int
type FloatFormatter func(float64) float64

const (
	TypeStringClear        = "clear"
	TypeStringTitle        = "title"
	TypeStringEveryTitle   = "every-title"
	TypeStringCode         = "code"
	TypeStringName         = "name"
	TypeStringAlpha        = "alpha"
	TypeStringNumeric      = "num"
	TypeStringAlphaNumeric = "alphanum"
	TypeStringLower        = "lower"
	TypeStringUpper        = "upper"
	TypeStringTrim         = "trim"
	TypeStringCyrillic     = "cyrillic"
)

const (
	TypeIntegerAbs = "abs"
)

const (
	TypeFloatAbs   = "abs"
	TypeFloatCeil  = "ceil"
	TypeFloatFloor = "floor"
)

var DefaultConverter = NewConverter().
	RegisterString(TypeStringClear, Clear).
	RegisterString(TypeStringTitle, Title).
	RegisterString(TypeStringEveryTitle, EveryTitle).
	RegisterString(TypeStringCode, Code).
	RegisterString(TypeStringName, Name).
	RegisterString(TypeStringAlpha, Alpha).
	RegisterString(TypeStringNumeric, Numeric).
	RegisterString(TypeStringAlphaNumeric, AlphaNumeric).
	RegisterString(TypeStringLower, strings.ToLower).
	RegisterString(TypeStringUpper, strings.ToUpper).
	RegisterString(TypeStringTrim, strings.TrimSpace).
	RegisterString(TypeStringCyrillic, Cyrillic).
	RegisterInteger(TypeIntegerAbs, Abs).
	RegisterFloat(TypeFloatAbs, math.Abs).
	RegisterFloat(TypeFloatCeil, math.Ceil).
	RegisterFloat(TypeFloatFloor, math.Floor)

type Converter struct {
	stringFormatters map[string]StringFormatter
	intFormatters    map[string]IntFormatter
	floatFormatters  map[string]FloatFormatter
}

func NewConverter() *Converter {
	return &Converter{
		stringFormatters: make(map[string]StringFormatter),
		intFormatters:    make(map[string]IntFormatter),
		floatFormatters:  make(map[string]FloatFormatter),
	}
}

func (converter *Converter) RegisterString(name string, formatter StringFormatter) *Converter {
	converter.stringFormatters[name] = formatter
	return converter
}

func (converter *Converter) RegisterInteger(name string, formatter IntFormatter) *Converter {
	converter.intFormatters[name] = formatter
	return converter
}

func (converter *Converter) RegisterFloat(name string, formatter FloatFormatter) *Converter {
	converter.floatFormatters[name] = formatter
	return converter
}

func (converter *Converter) Convert(input any) (err error) {
	defer errs.Wrap("Format", &err, "Convert")
	defer func() {
		if r := recover(); r != nil {
			err = try.CatchPanic(r)
		}
	}()

	// check input object for pointer type
	t := flex.Type(input)
	if !t.IsPtr() {
		return errors.New("input is not a pointer")
	}

	// prepare reflect objects
	inputType := t.Unwrap().Type()
	value := reflect.ValueOf(input).Elem()

	// check every field of provided object
	for i := 0; i < inputType.NumField(); i++ {
		// if field has no "format" tag - skip
		tag, ok := inputType.Field(i).Tag.Lookup("format")
		if !ok {
			continue
		}

		// get tag value
		field := value.Field(i)

		// update field value
		formatted, wasFormat := converter.format(inputType.Field(i).Type, field, tag)
		if !wasFormat {
			continue
		}

		field.Set(formatted)
	}

	return nil
}

func (converter *Converter) format(inputType reflect.Type, field reflect.Value, tag string) (formatted reflect.Value, ok bool) {
	switch inputType.Kind() {
	case reflect.String:
		return converter.formatString(field, tag), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return converter.formatInt(inputType.Kind(), field, tag), true
	case reflect.Float32, reflect.Float64:
		return converter.formatFloat(inputType.Kind(), field, tag), true
	default:
		return formatted, false
	}
}

func (converter *Converter) formatString(field reflect.Value, tag string) reflect.Value {
	// get tag value
	fieldValue := field.String()

	// format field value by provided tags
	tags := strings.Split(tag, ",")
	for _, fieldTag := range tags {
		// if tag value not found - skip
		formatter, formatterMatch := converter.stringFormatters[fieldTag]
		if !formatterMatch {
			continue
		}

		// format field value
		fieldValue = formatter(fieldValue)
	}

	return reflect.ValueOf(fieldValue)
}

func (converter *Converter) formatInt(kind reflect.Kind, field reflect.Value, tag string) reflect.Value {
	// get tag value
	fieldValue := int(field.Int())

	// format field value by provided tags
	tags := strings.Split(tag, ",")
	for _, fieldTag := range tags {
		// if tag value not found - skip
		formatter, formatterMatch := converter.intFormatters[fieldTag]
		if !formatterMatch {
			continue
		}

		// format field value
		fieldValue = formatter(fieldValue)
	}

	switch kind {
	case reflect.Int:
		return reflect.ValueOf(fieldValue)
	case reflect.Int8:
		return reflect.ValueOf(int8(fieldValue))
	case reflect.Int16:
		return reflect.ValueOf(int16(fieldValue))
	case reflect.Int32:
		return reflect.ValueOf(int32(fieldValue))
	case reflect.Int64:
		return reflect.ValueOf(int64(fieldValue))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(fieldValue))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(fieldValue))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(fieldValue))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(fieldValue))
	default:
		return reflect.ValueOf(fieldValue)
	}
}

func (converter *Converter) formatFloat(kind reflect.Kind, field reflect.Value, tag string) reflect.Value {
	// get tag value
	fieldValue := field.Float()

	// format field value by provided tags
	tags := strings.Split(tag, ",")
	for _, fieldTag := range tags {
		// if tag value not found - skip
		formatter, formatterMatch := converter.floatFormatters[fieldTag]
		if !formatterMatch {
			continue
		}

		// format field value
		fieldValue = formatter(fieldValue)
	}

	switch kind {
	case reflect.Float32:
		return reflect.ValueOf(float32(fieldValue))
	case reflect.Float64:
		return reflect.ValueOf(fieldValue)
	default:
		return reflect.ValueOf(fieldValue)
	}
}

func Convert(input any) error {
	return DefaultConverter.Convert(input)
}
