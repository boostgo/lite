package flex

import (
	"fmt"
	"github.com/boostgo/lite/collections/list"
	"reflect"
	"strings"
)

type ObjectField interface {
	Type() reflect.Type
	IsStruct() bool
	Fields() []ObjectField
	Tag(tagName string) []string

	fmt.Stringer
}

type objectField struct {
	t ObjectType
	f reflect.StructField
}

func Field(field reflect.StructField) ObjectField {
	return newObjectField(field)
}

func newObjectField(field reflect.StructField) *objectField {
	return &objectField{
		f: field,
		t: Type(field.Type),
	}
}

func (f *objectField) Type() reflect.Type {
	return f.f.Type
}

func (f *objectField) IsStruct() bool {
	return f.t.IsStruct()
}

func (f *objectField) Fields() []ObjectField {
	str, err := Struct(reflect.New(f.t.Type()).Elem().Interface())
	if err != nil {
		return []ObjectField{}
	}

	return str.Fields()
}

func (f *objectField) Tag(tagName string) []string {
	tagValue := f.f.Tag.Get(tagName)
	if tagValue == "" {
		return nil
	}

	return list.Of(strings.Split(tagValue, ",")).Map(strings.TrimSpace)
}

func (f *objectField) String() string {
	return f.f.Name
}
