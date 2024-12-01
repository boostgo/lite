package format

import (
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/types/flex"
	"reflect"
	"strings"
)

type Formatter func(string) string

const (
	TypeTitle    = "title"
	TypeUsername = "username"
	TypeAlpha    = "alpha"
	TypeNumeric  = "num"
)

var DefaultConverter = NewConverter().
	Register(TypeTitle, Title).
	Register(TypeUsername, Username).
	Register(TypeAlpha, Alpha).
	Register(TypeNumeric, Numeric)

type Converter struct {
	formatters map[string]Formatter
}

func NewConverter() *Converter {
	return &Converter{
		formatters: make(map[string]Formatter),
	}
}

func (converter *Converter) Register(name string, formatter Formatter) *Converter {
	converter.formatters[name] = formatter
	return converter
}

func (converter *Converter) Convert(input any) (err error) {
	defer errs.Wrap("Format", &err, "Convert")

	t := flex.Type(input)
	if !t.IsPtr() {
		return errors.New("input is not a pointer")
	}

	inputType := t.Unwrap().Type()
	value := reflect.ValueOf(input).Elem()

	for i := 0; i < inputType.NumField(); i++ {
		tag, ok := inputType.Field(i).Tag.Lookup("format")
		if !ok {
			continue
		}

		if inputType.Field(i).Type.Kind() != reflect.String {
			continue
		}

		field := value.Field(i)
		fieldValue := field.String()

		tags := strings.Split(tag, ",")
		for _, fieldTag := range tags {
			formatter, formatterMatch := converter.formatters[fieldTag]
			if !formatterMatch {
				continue
			}

			fieldValue = formatter(fieldValue)
		}

		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}

func Convert(input any) (err error) {
	defer errs.Wrap("Format", &err, "Convert")

	t := flex.Type(input)
	if !t.IsPtr() {
		return errors.New("input is not a pointer")
	}

	inputType := t.Unwrap().Type()
	value := reflect.ValueOf(input).Elem()

	for i := 0; i < inputType.NumField(); i++ {
		tag, ok := inputType.Field(i).Tag.Lookup("format")
		if !ok {
			continue
		}

		if inputType.Field(i).Type.Kind() != reflect.String {
			continue
		}

		field := value.Field(i)
		fieldValue := field.String()

		tags := strings.Split(tag, ",")
		for _, fieldTag := range tags {
			formatter, formatterMatch := DefaultConverter.formatters[fieldTag]
			if !formatterMatch {
				continue
			}

			fieldValue = formatter(fieldValue)
		}

		field.Set(reflect.ValueOf(fieldValue))
	}

	return nil
}
