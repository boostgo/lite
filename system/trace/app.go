package trace

import "sync/atomic"

var (
	_masterMode = atomic.Bool{}
)

func IAmMaster() {
	_masterMode.Store(true)
}

func AmIMaster() bool {
	return _masterMode.Load()
}
