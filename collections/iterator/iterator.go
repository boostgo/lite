package iterator

import (
	"errors"
	"github.com/boostgo/lite/types/param"
	"github.com/boostgo/lite/types/to"
)

type Iterator[T any] struct {
	current int
	params  []T
	err     error
}

func New[T any](params []T) *Iterator[T] {
	return &Iterator[T]{
		params: params,
	}
}

func (it *Iterator[T]) SetError(err error) *Iterator[T] {
	it.err = err
	return it
}

func (it *Iterator[T]) Skip(count int) *Iterator[T] {
	for i := 0; i < count; i++ {
		_, _ = it.Next()
	}
	return it
}

func (it *Iterator[T]) Next() (param.Param, bool) {
	if it.current >= len(it.params) {
		return param.Param{}, false
	}

	defer func() {
		it.current++
	}()
	return param.New(to.String(it.params[it.current])), true
}

func (it *Iterator[T]) NextString() (string, error) {
	nextParam, ok := it.Next()
	if !ok {
		if it.err != nil {
			return "", it.err
		}

		return "", errors.New("wrong params quantity error")
	}

	return nextParam.String(), nil
}

func (it *Iterator[T]) NextInt() (int, error) {
	nextParam, ok := it.Next()
	if !ok {
		if it.err != nil {
			return 0, it.err
		}

		return 0, errors.New("wrong params quantity error")
	}

	return nextParam.Int()
}
