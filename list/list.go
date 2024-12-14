package list

import (
	"fmt"
	"github.com/boostgo/lite/types/to"
	"math/rand"
)

type OfSlice[T any] interface {
	All(fn func(T) bool) bool
	Any(fn func(T) bool) bool
	Each(fn func(int, T)) OfSlice[T]
	EachErr(fn func(int, T) error) error
	Equal(against []T, fn func(T, T) bool) bool

	Filter(fn func(T) bool) OfSlice[T]
	FilterNot(fn func(T) bool) OfSlice[T]
	Reverse() OfSlice[T]
	Shuffle(source ...rand.Source) OfSlice[T]
	Sort(less func(a, b T) bool) OfSlice[T]
	Unique(fn func(a, b T) bool) []T

	Join(joins ...[]T) OfSlice[T]
	Add(elements ...T) OfSlice[T]
	AddLeft(elements ...T) OfSlice[T]
	Set(index int, elements ...T) OfSlice[T]
	Remove(index ...int) OfSlice[T]
	RemoveWhere(fn func(T) bool) OfSlice[T]
	IndexOf(fn func(T) bool) int
	Clear(capacity ...int) OfSlice[T]

	Single(fn func(T) bool) (T, bool)
	Exist(fn func(T) bool) bool
	First(fn func(T) bool) (T, bool)
	Last(fn func(T) bool) (T, bool)
	Contains(value T, fn ...func(T, T) bool) bool
	Get(index int) *T
	Slice() []T
	SliceAny(fn ...func(T) any) []any
	Sub(start, end int) OfSlice[T]
	Map(fn func(T) T) []T
	MapErr(fn func(T) (T, error)) ([]T, error)
	JoinString(fn func(T) string, sep ...string) string

	Len() int
	Cap() int
	fmt.Stringer
}

func Of[T any](list []T) OfSlice[T] {
	return newOf[T](list)
}

type ofSlice[T any] struct {
	source []T
}

func newOf[T any](source []T) *ofSlice[T] {
	return &ofSlice[T]{
		source: source,
	}
}

func (os *ofSlice[T]) All(fn func(T) bool) bool {
	return All(os.source, fn)
}

func (os *ofSlice[T]) Any(fn func(T) bool) bool {
	return Any(os.source, fn)
}

func (os *ofSlice[T]) Filter(fn func(T) bool) OfSlice[T] {
	return newOf(Filter(os.source, fn))
}

func (os *ofSlice[T]) FilterNot(fn func(T) bool) OfSlice[T] {
	return newOf(FilterNot(os.source, fn))
}

func (os *ofSlice[T]) Each(fn func(int, T)) OfSlice[T] {
	Each(os.source, fn)
	return os
}

func (os *ofSlice[T]) EachErr(fn func(int, T) error) error {
	return EachErr(os.source, fn)
}

func (os *ofSlice[T]) Equal(against []T, fn func(T, T) bool) bool {
	return AreEqual(os.source, against, fn)
}

func (os *ofSlice[T]) Single(fn func(T) bool) (T, bool) {
	return Single(os.source, fn)
}

func (os *ofSlice[T]) Exist(fn func(T) bool) bool {
	return Exist(os.source, fn)
}

func (os *ofSlice[T]) First(fn func(T) bool) (T, bool) {
	return First(os.source, fn)
}

func (os *ofSlice[T]) Last(fn func(T) bool) (T, bool) {
	return Last(os.source, fn)
}

func (os *ofSlice[T]) Contains(value T, fn ...func(T, T) bool) bool {
	return Contains(os.source, value, fn...)
}

func (os *ofSlice[T]) Get(index int) *T {
	return Get(os.source, index)
}

func (os *ofSlice[T]) Reverse() OfSlice[T] {
	return newOf(Reverse(os.source))
}

func (os *ofSlice[T]) Shuffle(source ...rand.Source) OfSlice[T] {
	return newOf(Shuffle(os.source, source...))
}

func (os *ofSlice[T]) Sort(less func(a, b T) bool) OfSlice[T] {
	return newOf(Sort(os.source, less))
}

func (os *ofSlice[T]) Add(elements ...T) OfSlice[T] {
	return newOf(Add(os.source, elements...))
}

func (os *ofSlice[T]) Join(joins ...[]T) OfSlice[T] {
	joins = append(joins, os.source)
	return newOf(Join(joins...))
}

func (os *ofSlice[T]) AddLeft(elements ...T) OfSlice[T] {
	return newOf(AddLeft(os.source, elements...))
}

func (os *ofSlice[T]) Set(index int, elements ...T) OfSlice[T] {
	return newOf(Set(os.source, index, elements...))
}

func (os *ofSlice[T]) Remove(index ...int) OfSlice[T] {
	return newOf(Remove(os.source, index...))
}

func (os *ofSlice[T]) IndexOf(fn func(T) bool) int {
	return IndexOf(os.source, fn)
}

func (os *ofSlice[T]) RemoveWhere(fn func(T) bool) OfSlice[T] {
	return newOf(RemoveWhere(os.source, fn))
}

func (os *ofSlice[T]) Clear(capacity ...int) OfSlice[T] {
	newCapacity := 0
	if len(capacity) > 0 {
		newCapacity = capacity[0]
	}

	os.source = make([]T, 0, newCapacity)
	return os
}

func (os *ofSlice[T]) Slice() []T {
	return os.source
}

func (os *ofSlice[T]) SliceAny(fn ...func(T) any) []any {
	return SliceAny(os.source, fn...)
}

func (os *ofSlice[T]) Sub(start, end int) OfSlice[T] {
	return newOf(Sub(os.source, start, end))
}

func (os *ofSlice[T]) Map(fn func(T) T) []T {
	return Map(os.source, fn)
}

func (os *ofSlice[T]) MapErr(fn func(T) (T, error)) ([]T, error) {
	return MapErr(os.source, fn)
}

func (os *ofSlice[T]) JoinString(fn func(T) string, sep ...string) string {
	return JoinString(os.source, fn, sep...)
}

func (os *ofSlice[T]) Len() int {
	return len(os.source)
}

func (os *ofSlice[T]) Cap() int {
	return cap(os.source)
}

func (os *ofSlice[T]) String() string {
	return to.String(os.source)
}

func (os *ofSlice[T]) Unique(fn func(a, b T) bool) []T {
	return Unique(os.source, fn)
}
