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

// Error is custom error which implements built-in error interface.
// Struct contains hierarchy of error messages and their types; context (map) and inner error.
// For example, error types could be like "User Handler - User Usecase - User Repository - SQL"
// It means that first error created on "SQL" level (sql, sqlx or any other module), then error wrapped
// by "User Repository" level, then "User Usecase" level and so on.
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

// Copy copies provided err to the new one.
// Inner errors sets inside new error as one inner error.
// If inner errors contains only 1 error it will be 1 error, if errors more than 1, it will be "Join error"
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

// Copy copies current error to the new one.
// Inner errors sets inside new error as one inner error.
// If inner errors contains only 1 error it will be 1 error, if errors more than 1, it will be "Join error"
func (err *Error) Copy(innerErrors ...error) error {
	return Copy(err, innerErrors...)
}

// Message returns all messages joined in one and reverse them.
// For example, there are messages: ["QueryxContext", "GetUser", "GetByID"]
// it will be "QueryxContext - GetByID - GetUser"
func (err *Error) Message() string {
	return strings.Join(list.Reverse(err.message), " - ")
}

// SetType append new type in chain of errors
func (err *Error) SetType(errorType string) *Error {
	err.errorTypes = append(err.errorTypes, errorType)
	return err
}

// Type returns all types joined in one and reverse them.
// For example, there are types: ["SQL", "User Repository", "User Usecase"]
// it will be "User Usecase - User Repository - SQL"
func (err *Error) Type() string {
	return strings.Join(list.Reverse(err.errorTypes), " - ")
}

// Context returns current error context (map)
func (err *Error) Context() map[string]any {
	return err.context
}

// SetContext append all key-value pairs to the current context map
func (err *Error) SetContext(context map[string]any) *Error {
	if context == nil || len(context) == 0 {
		return err
	}

	for key, value := range context {
		err.context[key] = value
	}

	return err
}

// AddContext append new key-value one pair to the current context map.
// But if provided key-value pair is array string as value and "trace" as key, it will be ignored
func (err *Error) AddContext(key string, value any) *Error {
	if value == nil {
		return err
	}

	if arr, ok := value.([]string); ok && key == "trace" {
		if len(arr) == 0 {
			return err
		}
	}

	err.context[key] = value

	return err
}

// InnerError returns inner error
func (err *Error) InnerError() error {
	return err.innerError
}

// SetError sets inner error.
// If inner errors more than 1 it will be "join error", if error is 1 it will be provided by itself
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

// Error returns result of String() method
func (err *Error) Error() string {
	return err.String()
}

// String returns string representation of current error.
// Method uses string builder and it's grow method.
// Method prints: types, messages and context
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

// Is compares current error with provided target error.
// By comparing errors method check if provided error is custom or not,
// if it is custom - use equals method.
// If it is not custom - unwrap current error and compare unwrapped inner errors with provided target
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

// Unwrap takes inner error and try to take inside wrapped errors.
// Method works only for custom errors, otherwise to result error slice will be added just inner error by itself
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

// setMessage appends new message to the chain
func (err *Error) setMessage(message string) *Error {
	err.message = append(err.message, message)
	return err
}

// grow - calculates length of message going to be built
func (err *Error) grow() int {
	var grow int
	// count error types + space and "-" symbol
	if len(err.errorTypes) > 0 {
		for i := 0; i < len(err.errorTypes); i++ {
			grow += len(err.errorTypes[i]) + 2
		}
	}

	// count inner error text
	if err.innerError != nil {
		grow += len(err.innerError.Error()) + 2
	}

	// count all context vars' lengths
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

// equals compare two provided custom errors by fields like:
// "type" and "error string"
func equals(err, target *Error) bool {
	return err.Type() == target.Type() &&
		err.Error() == target.Error()
}

// TryGet convert provided error to the custom and say it is custom or not
func TryGet(err error) (*Error, bool) {
	var custom *Error
	ok := errors.As(err, &custom)
	return custom, ok
}

// Get convert provided error to the custom and if it is not custom - return nil
func Get(err error) *Error {
	custom, ok := TryGet(err)
	if !ok {
		return nil
	}

	return custom
}

// IsType try to convert provided error to custom and compare error types
func IsType(err error, errorType string) bool {
	custom, ok := TryGet(err)
	if !ok {
		return false
	}

	return custom.Type() == errorType
}

// Is compare provided errors.
// Event if provided errors is not custom, comparing becomes by built-in "Is" function
func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	// if provided error is not custom, then compare by built in "Is" function
	errCustom, isCustom := TryGet(err)
	if !isCustom {
		return errors.Is(err, target)
	}

	// if provided target error is not custom, then compare by built in "Is" function
	targetCustom, isCustom := TryGet(target)
	if !isCustom {
		return errors.Is(err, target)
	}

	// if both errors are custom, compare by custom "Is" function
	return errCustom.Is(targetCustom)
}

// Wrap convert provided error to custom with the provided error type and message.
// If provided error is built-in (default), then it will be converted to custom.
// If it is already custom, just take custom and set to it one more type & message
func Wrap(errType string, err *error, message string, ctx ...map[string]any) {
	if *err != nil {
		var applyContext map[string]any
		if len(ctx) > 0 {
			applyContext = ctx[0]
		}

		custom, ok := TryGet(*err)
		if !ok {
			*err = New(message).
				SetType(errType).
				SetError(*err).
				SetContext(applyContext)
		} else {
			*err = custom.
				SetType(errType).
				setMessage(message).
				SetContext(applyContext)
		}
	}
}

// Type returns type of custom error.
// If provided error is built-in - return DefaultType
func Type(err error) string {
	custom, ok := TryGet(err)
	if !ok {
		return DefaultType
	}

	return custom.Type()
}
