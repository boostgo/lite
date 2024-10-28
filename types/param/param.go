package param

import (
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/types/to"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

func ErrorParseIntParam(err error, value string) error {
	return errs.
		New("Parse int param error").
		SetType("ParseIntParamError").
		SetError(err, errs.ErrUnprocessableEntity).
		AddContext("value", value)
}

func ErrorParseFloatParam(err error, value string) error {
	return errs.
		New("Parse float param error").
		SetType("ParseFloatParamError").
		SetError(err, errs.ErrUnprocessableEntity).
		AddContext("value", value)
}

func ErrorParseUUIDParam(err error, value string) error {
	return errs.
		New("Parse UUID param error").
		SetType("ParseUUIDParamError").
		SetError(err, errs.ErrUnprocessableEntity).
		AddContext("value", value)
}

type Param struct {
	value string
}

func New(value string) Param {
	return Param{
		value: value,
	}
}

func (param Param) IsEmpty() bool {
	return param.value == ""
}

func (param Param) String() string {
	return param.value
}

func (param Param) Strings() []string {
	return strings.Split(param.value, ",")
}

func (param Param) Int() (int, error) {
	intValue, err := strconv.Atoi(param.value)
	if err != nil {
		return 0, ErrorParseIntParam(err, param.value)
	}

	return intValue, nil
}

func (param Param) MustInt(defaultValue ...int) int {
	intValue, err := strconv.Atoi(param.value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return intValue
}

func (param Param) Int64() (int64, error) {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		return 0, ErrorParseIntParam(err, param.value)
	}

	return intValue, nil
}

func (param Param) MustInt64(defaultValue ...int64) int64 {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return intValue
}

func (param Param) Int32() (int32, error) {
	intValue, err := strconv.ParseInt(param.value, 10, 64)
	if err != nil {
		return 0, ErrorParseIntParam(err, param.value)
	}

	return int32(intValue), nil
}

func (param Param) MustInt32(defaultValue ...int32) int32 {
	intValue, err := strconv.ParseInt(param.value, 10, 32)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return int32(intValue)
}

func (param Param) Float32() (float32, error) {
	floatValue, err := strconv.ParseFloat(param.value, 32)
	if err != nil {
		return 0, ErrorParseFloatParam(err, param.value)
	}

	return float32(floatValue), nil
}

func (param Param) Float64() (float64, error) {
	floatValue, err := strconv.ParseFloat(param.value, 64)
	if err != nil {
		return 0, ErrorParseFloatParam(err, param.value)
	}

	return floatValue, nil
}

func (param Param) Bool() bool {
	return strings.ToLower(param.value) == "true"
}

func (param Param) UUID() (uuid.UUID, error) {
	uuidValue, err := uuid.Parse(param.value)
	if err != nil {
		return uuid.UUID{}, ErrorParseUUIDParam(err, param.value)
	}

	return uuidValue, nil
}

func (param Param) MustUUID() uuid.UUID {
	uuidValue, err := uuid.Parse(param.value)
	if err != nil {
		return uuid.UUID{}
	}

	return uuidValue
}

func (param Param) Bytes() []byte {
	return to.Bytes(param.value)
}
