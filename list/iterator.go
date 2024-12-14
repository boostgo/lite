package list

import "math/rand"

type Iterator[T any] interface {
	Next() bool
	Value() (T, bool)
	MustValue() T
	Skip(count int) Iterator[T]
	Reverse() Iterator[T]
	Shuffle(source ...rand.Source) Iterator[T]
	Each(fn func(int, T)) Iterator[T]
}

type iterator[T any] struct {
	source OfSlice[T]
	index  int
}

func Iterate[T any](source []T) Iterator[T] {
	return &iterator[T]{
		source: Of(source),
	}
}

func (it *iterator[T]) Reverse() Iterator[T] {
	it.source = it.source.Reverse()
	return it
}

func (it *iterator[T]) Shuffle(source ...rand.Source) Iterator[T] {
	it.source = it.source.Shuffle(source...)
	return it
}

func (it *iterator[T]) Next() bool {
	return it.index < it.source.Len()
}

func (it *iterator[T]) Value() (T, bool) {
	if it.Next() {
		item := *it.source.Get(it.index)
		it.index++
		return item, true
	}

	var zeroValue T
	return zeroValue, false
}

func (it *iterator[T]) MustValue() T {
	value, _ := it.Value()
	return value
}

func (it *iterator[T]) Skip(count int) Iterator[T] {
	it.index += count
	return it
}

func (it *iterator[T]) Each(fn func(int, T)) Iterator[T] {
	for it.Next() {
		idx := it.index

		value, ok := it.Value()
		if !ok {
			break
		}

		fn(idx, value)
	}
	return it
}
