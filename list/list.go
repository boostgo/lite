package list

import (
	"fmt"
	"github.com/boostgo/lite/types/to"
	"math/rand"
)

// OfSlice wrap over slice with helpful methods.
//
// Important: every method creates new slice but not modify current slice
type OfSlice[T any] interface {
	// All iterate over all slice elements and stop when func returns false.
	//
	// If was iterated all elements - returns true
	All(fn func(T) bool) bool
	// Any find element by provided condition func and if found - returns true
	Any(fn func(T) bool) bool
	// Each iterate over all slice elements and run provided function
	Each(fn func(int, T)) OfSlice[T]
	// EachErr iterate over all slice elements and run provided function.
	//
	// Stop iterating when provided function returns error
	EachErr(fn func(int, T) error) error
	// Equal compares with provided slice by using provided func
	Equal(against []T, fn func(T, T) bool) bool

	// Filter slice with provided condition func.
	//
	// Element appends to new slice if condition func returns true
	Filter(fn func(T) bool) OfSlice[T]
	// FilterNot slice with provided condition func.
	//
	// Element appends to new slice if condition func returns false
	FilterNot(fn func(T) bool) OfSlice[T]
	// Reverse slice
	Reverse() OfSlice[T]
	// Shuffle set elements in slice by random indexes.
	//
	// Could be provided custom rand.Source implementation
	Shuffle(source ...rand.Source) OfSlice[T]
	// Sort slice by provided compare function
	Sort(less func(a, b T) bool) OfSlice[T]
	// Unique make slice unique by provided condition func.
	Unique(fn func(a, b T) bool) []T

	// Join unions provided slices into one
	Join(joins ...[]T) OfSlice[T]
	// Add appends new elements to slice
	Add(elements ...T) OfSlice[T]
	// AddLeft append new elements to the start of slice
	AddLeft(elements ...T) OfSlice[T]
	// Set append new elements to slice on provided index
	Set(index int, elements ...T) OfSlice[T]
	// Remove delete elements from slice by provided indexes
	Remove(index ...int) OfSlice[T]
	// RemoveWhere delete element from slice by provided condition func
	RemoveWhere(fn func(T) bool) OfSlice[T]
	// IndexOf return index of found element by provided condition func.
	//
	// If element not found - returns -1
	IndexOf(fn func(T) bool) int
	// Clear recreate current slice with new capacity (which could be provided)
	Clear(capacity ...int) OfSlice[T]

	// Single returns element and found boolean by provided condition func.
	//
	// If element found - returns true
	Single(fn func(T) bool) (T, bool)
	// Exist check if element exist by provided condition func.
	//
	// The check performs on first matched element
	Exist(fn func(T) bool) bool
	// First returns element and found boolean by provided condition func.
	//
	// If element found - returns true.
	//
	// Element start matching from the start of the slice
	First(fn func(T) bool) (T, bool)
	// Last returns element and found boolean by provided condition func.
	//
	// If element found - returns true.
	//
	// Element start matching from the end of the slice
	Last(fn func(T) bool) (T, bool)
	// Contains check if element exist in slice.
	//
	// Could be provided custom comparing function, by default compares by using reflect.DeepEqual.
	//
	// The check performs on first matched element.
	Contains(value T, fn ...func(T, T) bool) bool
	// Get return element by index.
	//
	// If index is out of slice range - returns empty value of slice type
	Get(index int) T
	// Slice returns OfSlice slice
	Slice() []T
	// SliceAny convert slice to "any" type elements slice
	SliceAny(fn ...func(T) any) []any
	// Sub return "sub slice" by provided start & end indexes.
	Sub(start, end int) OfSlice[T]
	// Map create new slice and fill with values by provided convert func
	Map(fn func(T) T) []T
	// MapErr create new slice and fill with values by provided convert func and may return error
	MapErr(fn func(T) (T, error)) ([]T, error)
	// JoinString build string from slice elements.
	//
	// Every element string builds from provided func.
	//
	// Could be provided custom separator between element strings
	JoinString(fn func(T) string, sep ...string) string

	// Len returns length of slice
	Len() int
	// Cap returns capacity of slice
	Cap() int
	fmt.Stringer
}

// Of create OfSlice
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

func (os *ofSlice[T]) Get(index int) T {
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
