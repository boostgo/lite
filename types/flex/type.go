package flex

import (
	"reflect"

	"github.com/boostgo/collection/slicex"
)

type ObjectType interface {
	With(func(reflectType reflect.Type))
	Unwrap() ObjectType
	Is(kinds ...reflect.Kind) bool
	IsType(object any) bool

	IsPrimitive() bool
	IsString() bool
	IsBool() bool
	IsNumeric() bool

	IsReference() bool
	IsPtr() bool
	IsInterface() bool
	IsFunc() bool
	IsStruct(element ...any) bool
	IsSlice(element ...any) bool

	IsTime() bool
	IsUUID() bool
	IsBytes() bool

	String() string
	Type() reflect.Type
}

type objectType struct {
	t reflect.Type
}

func Type(object any) ObjectType {
	t, isType := object.(reflect.Type)
	if !isType {
		t = reflect.TypeOf(object)
	}
	return &objectType{
		t: t,
	}
}

func (t *objectType) With(with func(reflectType reflect.Type)) {
	with(t.t)
}

func (t *objectType) Unwrap() ObjectType {
	if !t.compare(reflect.Ptr) {
		return t
	}

	elem := t.t.Elem()
	for {
		if elem.Kind() != reflect.Ptr {
			break
		}

		elem = elem.Elem()
	}

	t.t = elem
	return t
}

func (t *objectType) Is(kinds ...reflect.Kind) bool {
	return t.compare(kinds...)
}

func (t *objectType) IsType(object any) bool {
	return Type(object).Unwrap().String() == t.String()
}

func (t *objectType) IsPtr() bool {
	return t.compare(reflect.Pointer)
}

func (t *objectType) IsInterface() bool {
	return t.compare(reflect.Interface)
}

func (t *objectType) IsString() bool {
	return t.compare(reflect.String)
}

func (t *objectType) IsBool() bool {
	return t.compare(reflect.Bool)
}

func (t *objectType) IsNumeric() bool {
	return t.compare(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128)
}

func (t *objectType) IsStruct(element ...any) bool {
	isStruct := t.compare(reflect.Struct)
	if !isStruct {
		return false
	}

	if len(element) > 0 {
		if element[0] == nil {
			return false
		}

		return t.t.String() == Type(element[0]).String()
	}

	return true
}

func (t *objectType) IsSlice(element ...any) bool {
	isSlice := t.compare(reflect.Slice, reflect.Array)
	if !isSlice {
		return false
	}

	if len(element) > 0 {
		if element[0] == nil {
			return false
		}

		return t.t.Elem().String() == Type(element[0]).Type().String()
	}

	return true
}

func (t *objectType) IsFunc() bool {
	return t.compare(reflect.Func)
}

func (t *objectType) IsPrimitive() bool {
	switch t.t.Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Interface,
		reflect.Slice, reflect.Array, reflect.Map:
		return false
	}

	return true
}

func (t *objectType) IsReference() bool {
	return !t.IsPrimitive()
}

func (t *objectType) IsBytes() bool {
	elem := byte(1)
	return t.IsSlice(elem)
}

func (t *objectType) IsTime() bool {
	return t.String() == "time.Time"
}

func (t *objectType) IsUUID() bool {
	return t.String() == "uuid.UUID"
}

func (t *objectType) String() string {
	return t.t.String()
}

func (t *objectType) Type() reflect.Type {
	return t.t
}

func (t *objectType) compare(destinations ...reflect.Kind) bool {
	return slicex.
		Of(destinations).
		Exist(func(kind reflect.Kind) bool {
			return t.t.Kind() == kind
		})
}
