package async

import (
	"sync"
)

// String representation of "thread safe" or "atomic" for string.
type String struct {
	value string
	mx    sync.RWMutex
}

func NewString(value string) *String {
	return &String{
		value: value,
	}
}

func (str *String) Load() string {
	str.mx.RLock()
	defer str.mx.RUnlock()
	return str.value
}

func (str *String) Store(value string) *String {
	str.mx.Lock()
	defer str.mx.Unlock()
	str.value = value
	return str
}

func (str *String) Append(appendValue string) *String {
	str.Store(str.value + appendValue)
	return str
}
