package errs

import (
	"github.com/boostgo/lite/types/to"
	"runtime/debug"
)

func Panic(err error) *Error {
	if custom, ok := TryGet(err); ok {
		return custom
	}

	return New("PANIC RECOVER").
		SetError(err).
		SetType("Panic").
		AddContext("trace", to.String(debug.Stack()))
}
