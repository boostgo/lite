package async

import (
	"github.com/boostgo/lite/list"
	"sync"
)

// Map thread safe map implementation.
//
// Contain data map defend by [sync.RWMutex].
//
// Getting keys list is cached (no need iteration).
//
// Can get length/size of [Map]
type Map[K comparable, V any] struct {
	keys list.OfSlice[K]
	data map[K]V
	mx   sync.RWMutex
}

// NewMap creates Map
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

// Store provided key & value pair
func (am *Map[K, V]) Store(key K, value V) *Map[K, V] {
	am.mx.Lock()
	defer am.mx.Unlock()

	am.data[key] = value
	am.keys = am.keys.Add(key)
	return am
}

// Load get value by provided key
func (am *Map[K, V]) Load(key K) (V, bool) {
	am.mx.RLock()
	defer am.mx.RUnlock()
	v, ok := am.data[key]
	return v, ok
}

// Keys return all keys
func (am *Map[K, V]) Keys() []K {
	am.mx.RLock()
	defer am.mx.RUnlock()
	return am.keys.Slice()
}

// Len returns length of map
func (am *Map[K, V]) Len() int {
	am.mx.RLock()
	defer am.mx.RUnlock()
	return len(am.data)
}

// Delete element by key
func (am *Map[K, V]) Delete(key K) *Map[K, V] {
	am.mx.Lock()
	defer am.mx.Unlock()

	delete(am.data, key)
	am.keys = am.keys.RemoveWhere(func(keyIterator K) bool {
		return key == keyIterator
	})

	return am
}

// Each iterate over all map elements.
// Stops when provided function return false or when all were keys iterated
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

// Map returns inner map
func (am *Map[K, V]) Map() map[K]V {
	return am.data
}
