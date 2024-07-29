package concurrent

import (
	"github.com/boostgo/lite/collections/list"
	"sync"
)

type Map[K comparable, V any] struct {
	keys list.OfSlice[K]
	data map[K]V
	mx   sync.RWMutex
}

func NewMap[K comparable, V any](size ...int) *Map[K, V] {
	mapSize := 8
	if len(size) > 0 {
		mapSize = size[0]
	}

	return &Map[K, V]{
		keys: list.Of([]K{}),
		data: make(map[K]V, mapSize),
	}
}

func (am *Map[K, V]) Store(key K, value V) *Map[K, V] {
	am.mx.Lock()
	defer am.mx.Unlock()

	am.data[key] = value
	am.keys = am.keys.Add(key)
	return am
}

func (am *Map[K, V]) Load(key K) (V, bool) {
	am.mx.RLock()
	defer am.mx.RUnlock()
	v, ok := am.data[key]
	return v, ok
}

func (am *Map[K, V]) Keys() []K {
	am.mx.RLock()
	defer am.mx.RUnlock()
	return am.keys.Slice()
}

func (am *Map[K, V]) Len() int {
	am.mx.RLock()
	defer am.mx.RUnlock()
	return len(am.data)
}

func (am *Map[K, V]) Delete(key K) *Map[K, V] {
	am.mx.Lock()
	defer am.mx.Unlock()

	delete(am.data, key)
	am.keys = am.keys.RemoveWhere(func(keyIterator K) bool {
		return key == keyIterator
	})

	return am
}

func (am *Map[K, V]) Each(fn func(key K, value V) bool) *Map[K, V] {
	for _, key := range am.Keys() {
		value, ok := am.Load(key)
		if !ok {
			continue
		}

		if !fn(key, value) {
			break
		}
	}

	return am
}