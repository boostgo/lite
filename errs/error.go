package errs

import (
	"errors"
	"github.com/boostgo/lite/collections/list"
	"github.com/boostgo/lite/types/to"
	"strings"
)

const (
	DefaultType = ""
)

type Error struct {
	message    []string
	errorTypes []string
	context    map[string]any
	innerError error
}

// New creates new Error object with provided message
func New(message string) *Error {
	messages := make([]string, 0)
	messages = append(messages, message)

	return &Error{
		message:    messages,
		errorTypes: make([]string, 0),
		context:    make(map[string]any),
	}
}

func Copy(err error, innerErrors ...error) error {
	custom, ok := TryGet(err)
	if !ok {
		return New(err.Error()).
			SetError(innerErrors...)
	}

	inner := make([]error, 0, len(innerErrors)+1)
	inner = append(inner, custom.innerError)
	inner = append(inner, innerErrors...)

	return New(custom.Message()).
		SetType(custom.Type()).
		SetContext(custom.Context()).
		SetError(inner...)
}

func (err *Error) Copy() error {
	return Copy(err)
}

func (err *Error) Message() string {
	return strings.Join(list.Reverse(err.message), " - ")
}

func (err *Error) SetType(errorType string) *Error {
	err.errorTypes = append(err.errorTypes, errorType)
	return err
}

func (err *Error) Type() string {
	return strings.Join(list.Reverse(err.errorTypes), " - ")
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
	builder.Grow(err.grow())

	if len(err.errorTypes) > 0 {
		builder.WriteString("[")
		builder.WriteString(err.Type())
		builder.WriteString("] ")
	}
	builder.WriteString(err.Message())

	if err.innerError != nil {
		innerMessage := err.innerError.Error()
		builder.WriteString(": ")
		builder.WriteString(innerMessage)
	}

	if err.context != nil && len(err.context) > 0 {
		builder.WriteString(". Context: ")
		for key, value := range err.context {
			if key == "trace" {
				trace, ok := value.([]string)
				if !ok {
					builder.WriteString("\n\t")
					builder.WriteString(value.(string))
					continue
				}

				for _, traceLine := range trace {
					builder.WriteString("\n\t")
					builder.WriteString(traceLine)
				}
				continue
			}

			builder.WriteString(key)
			builder.WriteString("=")
			builder.WriteString(to.String(value))
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

func (err *Error) setMessage(message string) *Error {
	err.message = append(err.message, message)
	return err
}

func (err *Error) grow() int {
	var grow int
	if len(err.errorTypes) > 0 {
		for i := 0; i < len(err.errorTypes); i++ {
			grow += len(err.errorTypes[i]) + 2
		}
	}

	if err.innerError != nil {
		grow += len(err.innerError.Error()) + 2
	}

	if err.context != nil && len(err.context) > 0 {
		grow += 11
		for key, value := range err.context {
			if key == "trace" {
				trace, ok := value.([]string)
				if !ok {
					continue
				}

				for _, traceLine := range trace {
					grow += len(traceLine) + 5
				}
				continue
			}

			grow += len(key) + len(to.String(value)) + 2
		}
	}
	return grow
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
		custom, ok := TryGet(*err)
		if !ok {
			*err = New(message).SetType(errType).SetError(*err)
		} else {
			*err = custom.
				SetType(errType).
				setMessage(message)
		}
	}
}

func Type(err error) string {
	custom, ok := TryGet(err)
	if !ok {
		return DefaultType
	}

	return custom.Type()
}
