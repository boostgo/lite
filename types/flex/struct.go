package flex

import (
	"errors"
	"fmt"
	"reflect"
)

type ObjectStruct interface {
	FieldsCount() int
	Fields() []ObjectField
	FieldsRow() []ObjectField
	FieldByType(t reflect.Type) any

	MethodsCount() int
	Methods() []ObjectFunc
	FieldValueByTag(tagName, tagValue string) any

	fmt.Stringer
}

type objectStruct struct {
	t ObjectType
	v reflect.Value
}

func Struct(object any) (ObjectStruct, error) {
	return newStruct(object)
}

func newStruct(object any) (*objectStruct, error) {
	t := Type(object)
	if !t.IsStruct() {
		// try to unwrap type & check
		if !t.Unwrap().IsStruct() {
			return nil, errors.New("object is not struct")
		}
	}

	return &objectStruct{
		t: t,
		v: reflect.ValueOf(object),
	}, nil
}

func (s *objectStruct) FieldsCount() int {
	return s.t.Type().NumField()
}

func (s *objectStruct) Fields() []ObjectField {
	fieldCount := s.FieldsCount()

	fields := make([]ObjectField, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		fields = append(fields, Field(s.t.Type().Field(i)))
	}

	return fields
}

func (s *objectStruct) FieldsRow() []ObjectField {
	fieldCount := s.FieldsCount()

	fields := make([]ObjectField, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field := s.t.Type().Field(i)

		t := Type(field.Type).Unwrap()
		if t.IsStruct() && !t.IsTime() {
			subStruct, _ := Struct(reflect.New(t.Type()).Elem().Interface())
			fields = append(fields, subStruct.FieldsRow()...)
			continue
		}

		fields = append(fields, Field(field))
	}

	return fields
}

func (s *objectStruct) FieldByType(t reflect.Type) any {
	fields := s.Fields()
	for index, field := range fields {
		if field.Type() == t {
			return s.v.Field(index).Interface()
		}
	}
	return nil
}

func (s *objectStruct) MethodsCount() int {
	// non-pointer type + pointer type
	return s.methodsCountNonPtr() + s.methodsCountPtr()
}

func (s *objectStruct) methodsCountNonPtr() int {
	return s.t.Type().NumMethod()
}

func (s *objectStruct) methodsCountPtr() int {
	return reflect.PtrTo(s.t.Type()).NumMethod()
}

func (s *objectStruct) Methods() []ObjectFunc {
	methodCountNonPtr := s.methodsCountNonPtr()
	methodCountPtr := s.methodsCountPtr()
	methods := make([]ObjectFunc, 0, methodCountNonPtr+methodCountPtr)
	for i := 0; i < methodCountNonPtr; i++ {
		methods = append(methods, fromMethod(s.t.Type().Method(i)))
	}
	ptrType := reflect.PtrTo(s.t.Type())
	for i := 0; i < methodCountPtr; i++ {
		methods = append(methods, fromMethod(ptrType.Method(i)))
	}
	return methods
}

func (s *objectStruct) FieldValueByTag(tagName, tagValue string) any {
	for i := 0; i < s.t.Type().NumField(); i++ {
		searchTagValue := s.t.Type().Field(i).Tag.Get(tagName)
		if len(searchTagValue) == 0 {
			continue
		}

		if searchTagValue == tagValue {
			return s.v.Field(i).Interface()
		}
	}

	// still not found? try go to parent objects
	for i := 0; i < s.t.Type().NumField(); i++ {
		kind := Type(s.t.Type().Field(i).Type)
		if kind.IsPrimitive() || kind.IsTime() || kind.IsUUID() {
			continue
		}

		str, _ := Struct(s.v.Field(i).Interface())
		searchValue := str.FieldValueByTag(tagName, tagValue)
		if searchValue != nil {
			return searchValue
		}
	}

	return nil
}

func (s *objectStruct) String() string {
	return s.t.String()
}
