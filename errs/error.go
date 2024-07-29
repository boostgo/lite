package errs

import (
	"errors"
	"github.com/boostgo/lite/types/to"
	"strings"
)

const (
	DefaultType = ""
)

type Error struct {
	message    string
	errorType  string
	context    map[string]any
	innerError error
}

// New creates new Error object with provided message
func New(message string) *Error {
	return &Error{
		message:   message,
		errorType: DefaultType,
		context:   make(map[string]any),
	}
}

func (err *Error) Message() string {
	return err.message
}

func (err *Error) SetType(errorType string) *Error {
	err.errorType = errorType
	return err
}

func (err *Error) Type() string {
	return err.errorType
}

func (err *Error) Context() map[string]any {
	return err.context
}

func (err *Error) SetContext(context map[string]any) *Error {
	for key, value := range context {
		err.context[key] = value
	}

	return err
}

func (err *Error) AddContext(key string, value any) *Error {
	if value == nil {
		return err
	}

	if arr, ok := value.([]string); ok {
		if len(arr) == 0 {
			return err
		}
	}

	err.context[key] = value

	return err
}

func (err *Error) InnerError() error {
	return err.innerError
}

func (err *Error) SetError(innerError ...error) *Error {
	if len(innerError) == 0 {
		return err
	}

	var inner error
	if len(innerError) == 1 {
		inner = innerError[0]
	} else {
		inner = Join(innerError...)
	}
	err.innerError = inner
	return err
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) String() string {
	builder := strings.Builder{}
	builder.Grow(len(err.message))
	if err.errorType != DefaultType {
		builder.Grow(len(err.errorType) + 2)
		builder.WriteString("[")
		builder.WriteString(err.errorType)
		builder.WriteString("] ")
	}
	builder.WriteString(err.message)

	if err.innerError != nil {
		innerMessage := err.innerError.Error()
		builder.Grow(len(innerMessage) + 2)
		builder.WriteString(": ")
		builder.WriteString(innerMessage)
	}

	if err.context != nil && len(err.context) > 0 {
		builder.Grow(11)
		builder.WriteString(". Context: ")
		for key, value := range err.context {
			if key == "trace" {
				trace := value.([]string)

				for _, traceLine := range trace {
					builder.Grow(len(traceLine) + 5)
					builder.WriteString("\n\t")
					builder.WriteString(traceLine)
				}
				continue
			}

			valueString := to.String(value)
			builder.Grow(len(key) + len(valueString) + 2)

			builder.WriteString(key)
			builder.WriteString("=")
			builder.WriteString(valueString)
			builder.WriteString(";")
		}
	}

	return builder.String()
}

func (err *Error) Is(target error) bool {
	custom, ok := TryGet(target)
	if !ok {
		if innerErrs := err.Unwrap(); innerErrs != nil && len(innerErrs) > 0 {
			for _, inner := range innerErrs {
				if errors.Is(inner, target) {
					return true
				}
			}
		}

		return false
	}

	return equals(err, custom)
}

func (err *Error) Unwrap() []error {
	if err.innerError == nil {
		return []error{}
	}

	unwrapped := make([]error, 0)
	unwrapped = append(unwrapped, err.innerError)
	custom, ok := TryGet(err.innerError)
	if ok {
		unwrapped = append(unwrapped, custom.Unwrap()...)
	}

	return unwrapped
}

func equals(err, target *Error) bool {
	return err.Type() == target.Type() &&
		err.Error() == target.Error()
}

func TryGet(err error) (*Error, bool) {
	var custom *Error
	ok := errors.As(err, &custom)
	return custom, ok
}

func Get(err error) *Error {
	custom, ok := TryGet(err)
	if !ok {
		return nil
	}

	return custom
}

func IsType(err error, errorType string) bool {
	custom, ok := TryGet(err)
	if !ok {
		return false
	}

	return custom.Type() == errorType
}

func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	errCustom, isCustom := TryGet(err)
	if !isCustom {
		return errors.Is(err, target)
	}

	targetCustom, isCustom := TryGet(target)
	if !isCustom {
		return errors.Is(err, target)
	}

	return errCustom.Is(targetCustom)
}

func Wrap(errType string, err *error, message string) {
	if *err != nil {
		*err = New(message).SetType(errType).SetError(*err)
	}
}

func Type(err error) string {
	custom, ok := TryGet(err)
	if !ok {
		return DefaultType
	}

	return custom.Type()
}
