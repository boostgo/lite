package flex

import (
	"errors"
	"fmt"
	"reflect"
)

type ObjectSlice interface {
	Elements() []any

	fmt.Stringer
}

func Slice(slice any) (ObjectSlice, error) {
	return newObjectSlice(slice)
}

type objectSlice struct {
	t ObjectType
	v reflect.Value
}

func newObjectSlice(slice any) (*objectSlice, error) {
	t := Type(slice)
	if !t.IsSlice() {
		return nil, errors.New("object is not a slice")
	}

	return &objectSlice{
		t: t,
		v: reflect.ValueOf(slice),
	}, nil
}

func (s *objectSlice) Elements() []any {
	masked := make([]any, s.v.Len())
	for i := 0; i < s.v.Len(); i++ {
		masked[i] = s.v.Index(i).Interface()
	}
	return masked
}

func (s *objectSlice) String() string {
	return ""
}
