package list

import "math/rand"

// Iterator tool for iterating by provided slice.
//
// Use Next() and Value() methods to iterate & get values
type Iterator[T any] interface {
	// Next check is iterator reached last element
	Next() bool
	// Value return current iterator value & increment current index
	Value() (T, bool)
	// MustValue calls Value, but with no "exist" boolean
	MustValue() T
	// Skip skips iterator values for provided count
	Skip(count int) Iterator[T]
	// Reverse reverses slice inside Iterator
	Reverse() Iterator[T]
	// Shuffle shuffles slice inside Iterator
	Shuffle(source ...rand.Source) Iterator[T]
	// Each iterate every element of Iterator and stops when Value returns false
	Each(fn func(int, T)) Iterator[T]
}

type iterator[T any] struct {
	source OfSlice[T]
	index  int
}

// Iterate creates Iterator with provided slice
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
		item := it.source.Get(it.index)
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
