package set

import (
	"github.com/boostgo/lite/collections/list"
	"math/rand"
	"reflect"
	"sync"
)

type Set[T any] struct {
	data []T
	mx   sync.RWMutex
}

func New[T any](length ...int) *Set[T] {
	var capacity int
	if len(length) > 0 && length[0] > 0 {
		capacity = length[0]
	}

	return &Set[T]{
		data: make([]T, 0, capacity),
	}
}

func (set *Set[T]) Clear() *Set[T] {
	set.mx.Lock()
	defer set.mx.Unlock()
	capacity := cap(set.data)
	set.data = make([]T, 0, capacity)
	return set
}

func (set *Set[T]) Pop() T {
	set.mx.Lock()
	defer set.mx.Unlock()
	index := len(set.data) - 1
	defer set.remove(index)
	return set.data[index]
}

func (set *Set[T]) Add(item T) *Set[T] {
	set.mx.Lock()
	defer set.mx.Unlock()

	if set.exist(item) {
		return set
	}

	set.data = append(set.data, item)
	return set
}

func (set *Set[T]) Remove(index int) *Set[T] {
	set.mx.Lock()
	defer set.mx.Unlock()

	if len(set.data) == 0 {
		return set
	}

	return set.remove(index)
}

func (set *Set[T]) Shuffle() *Set[T] {
	set.mx.Lock()
	defer set.mx.Unlock()

	for i := range set.data {
		j := rand.Intn(i + 1)
		set.data[i], set.data[j] = set.data[j], set.data[i]
	}

	return set
}

func (set *Set[T]) Slice() []T {
	set.mx.RLock()
	defer set.mx.RUnlock()
	return set.data
}

func (set *Set[T]) Length() int {
	set.mx.RLock()
	defer set.mx.RUnlock()
	return len(set.data)
}

func (set *Set[T]) exist(searchItem T) bool {
	for _, item := range set.data {
		if reflect.DeepEqual(item, searchItem) {
			return true
		}
	}

	return false
}

func (set *Set[T]) remove(index int) *Set[T] {
	set.data = list.Remove(set.data, index)
	return set
}
