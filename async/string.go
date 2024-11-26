package async

import (
	"sync"
)

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

func (str *String) Store(value string) {
	str.mx.Lock()
	defer str.mx.Unlock()
	str.value = value
}
