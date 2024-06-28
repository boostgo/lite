package errs

import (
	"encoding/json"
	"errors"
	"github.com/boostgo/lite/types/content"
	"github.com/boostgo/lite/types/to"
	"net/http"
	"strings"
)

type Error struct {
	message    string
	errorType  *string
	httpCode   int
	context    map[string]any
	innerError error
}

type outputError struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Type    *string        `json:"type,omitempty"`
	Code    int            `json:"code"`
	Context map[string]any `json:"context,omitempty"`
}

// New creates new Boost Error object with given message
func New(message string) *Error {
	return &Error{
		message:  message,
		httpCode: http.StatusInternalServerError,
		context:  make(map[string]any),
	}
}

const (
	status = "ERROR"
)

func (err *Error) Message() string {
	return err.message
}

func (err *Error) SetHttpCode(code int) *Error {
	err.httpCode = code
	return err
}

func (err *Error) HttpCode() int {
	return err.httpCode
}

func (err *Error) SetType(errorType string) *Error {
	err.errorType = &errorType
	return err
}

func (err *Error) Type() *string {
	return err.errorType
}

func (err *Error) ContentType() string {
	return content.JSON
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

func (err *Error) SetError(innerError error) *Error {
	err.innerError = innerError
	return err
}

func (err *Error) Error() string {
	return err.String()
}

func (err *Error) JSON() []byte {
	errorMessage := err.message
	if err.innerError != nil {
		errorMessage += " | " + err.innerError.Error()
	}

	output := outputError{
		Status:  status,
		Message: errorMessage,
		Type:    err.errorType,
		Code:    err.httpCode,
		Context: err.context,
	}

	outputInBytes, _ := json.Marshal(output)
	return outputInBytes
}

func (err *Error) String() string {
	builder := strings.Builder{}
	builder.Grow(500)
	if err.errorType != nil {
		builder.WriteString("[")
		builder.WriteString(*err.errorType)
		builder.WriteString("] ")
	}
	builder.WriteString(err.message)

	if err.innerError != nil {
		builder.WriteString(": ")
		builder.WriteString(err.innerError.Error())
	}

	if err.context != nil && len(err.context) > 0 {
		builder.WriteString(". Context: ")
		for key, value := range err.context {
			if key == "trace" {
				trace := value.([]string)

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

func equals(err, target *Error) bool {
	return err.HttpCode() == target.HttpCode() &&
		err.Type() == target.Type() &&
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

	customErrType := custom.Type()
	return customErrType != nil && *customErrType == errorType
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

func Join(errors ...error) error {
	return newJoin(errors...)
}

func Wrap(errType string, err error) error {
	custom, ok := TryGet(err)
	if !ok {
		custom = New(err.Error()).SetType(errType)
	} else {
		_ = custom.SetType(errType)
	}
	return custom
}

func Type(err error) *string {
	custom, ok := TryGet(err)
	if !ok {
		return nil
	}

	return custom.Type()
}

func HttpCode(err error) int {
	custom, ok := TryGet(err)
	if !ok {
		return 0
	}

	return custom.httpCode
}

func FromBytes(response []byte) (*Error, bool) {
	var output outputError
	if err := json.Unmarshal(response, &output); err != nil {
		return nil, false
	}

	return &Error{
		message:   output.Message,
		errorType: output.Type,
		httpCode:  output.Code,
		context:   output.Context,
	}, true
}
