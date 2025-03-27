package flex

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/boostgo/collection/slicex"
)

type ObjectFunc interface {
	ArgumentsCount() int
	Arguments() []reflect.Type
	ArgumentsUnwrapped() []reflect.Type

	ReturnsCount() int
	Returns() []reflect.Type
	Call(arguments ...any) ([]any, error)
}

func Func(object any) (ObjectFunc, error) {
	return newFunc(object)
}

func fromMethod(method reflect.Method) ObjectFunc {
	return newFuncFromMethod(method)
}

type objectFunc struct {
	t ObjectType
	v reflect.Value
	m *reflect.Method
}

func newFunc(object any) (*objectFunc, error) {
	t := Type(object)
	if !t.IsFunc() {
		return nil, errors.New("object is not func")
	}

	return &objectFunc{
		t: t,
		v: reflect.ValueOf(object),
	}, nil
}

func newFuncFromMethod(method reflect.Method) *objectFunc {
	return &objectFunc{
		t: Type(method.Type),
		m: &method,
	}
}

func (f *objectFunc) isArgumentsAnySlice() bool {
	if f.ArgumentsCount() != 1 {
		return false
	}
	firstArg := f.funcType().In(0)
	if firstArg.Kind() != reflect.Slice {
		return false
	}
	return firstArg.Elem().Kind() == reflect.Interface
}

func (f *objectFunc) ArgumentsCount() int {
	return f.funcType().NumIn()
}

func (f *objectFunc) Arguments() []reflect.Type {
	argCount := f.ArgumentsCount()
	arguments := make([]reflect.Type, 0, argCount)
	for i := 0; i < argCount; i++ {
		arguments = append(arguments, f.funcType().In(i))
	}
	return arguments
}

func (f *objectFunc) ArgumentsUnwrapped() []reflect.Type {
	argCount := f.ArgumentsCount()
	arguments := make([]reflect.Type, 0, argCount)
	for i := 0; i < argCount; i++ {
		arguments = append(arguments, Type(f.funcType().In(i)).Unwrap().Type())
	}
	return arguments
}

func (f *objectFunc) ReturnsCount() int {
	return f.funcType().NumOut()
}

func (f *objectFunc) Returns() []reflect.Type {
	returnCount := f.ReturnsCount()
	returns := make([]reflect.Type, 0, returnCount)
	for i := 0; i < returnCount; i++ {
		returns = append(returns, f.funcType().Out(i))
	}
	return returns
}

func catchPanic(err any) error {
	if err == nil {
		return nil
	}

	fromAny := func(err any) string {
		if err == nil {
			return ""
		}

		switch v := err.(type) {
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		case error:
			return v.Error()
		}

		return ""
	}

	return errors.New(fromAny(err))
}

func (f *objectFunc) Call(arguments ...any) (returnValues []any, err error) {
	defer func() {
		if err == nil {
			err = catchPanic(recover())
		}
	}()

	argumentTypes := f.Arguments()
	isAnySlice := f.isArgumentsAnySlice()

	// check arguments count
	if !isAnySlice && len(arguments) != len(argumentTypes) {
		err = errors.New("arguments count does not match")
		return
	}

	// check arguments types
	for index, argType := range argumentTypes {
		if !isAnySlice && reflect.TypeOf(arguments[index]) != argType {
			err = errors.New("arguments types does not match")
			return
		}
	}

	// prepare in arguments
	in := make([]reflect.Value, 0, len(arguments))
	for _, arg := range arguments {
		in = append(in, reflect.ValueOf(arg))
	}

	// call function
	values := f.funcValue().Call(in)

	// convert result to any
	return slicex.Map(values, func(v reflect.Value) any {
		return v.Interface()
	}), nil
}

func (f *objectFunc) funcType() reflect.Type {
	if f.m != nil {
		return f.m.Type
	}
	return f.t.Type()
}

func (f *objectFunc) funcValue() reflect.Value {
	if f.m != nil {
		return f.m.Func
	}
	return f.v
}
