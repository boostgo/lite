package kafka

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/boostgo/lite/internal/reflectx"

	"github.com/IBM/sarama"
	"github.com/boostgo/convert"
	"github.com/boostgo/lite/system/validator"
)

var (
	_validator     *validator.Validator
	_validatorOnce sync.Once
)

func init() {
	_validatorOnce.Do(func() {
		_validator, _ = validator.New()
	})
}

// Parse message body to provided export object (which must be ptr) and validate for "validate" tags.
func Parse(message *sarama.ConsumerMessage, export any) error {
	// check export type
	if !reflectx.IsPointer(export) {
		return errors.New("export object must be pointer")
	}

	// parse message
	if err := json.Unmarshal(message.Value, export); err != nil {
		return err
	}

	// validate parsed message body
	return _validator.Struct(export)
}

// Header search header in provided message by header name.
func Header(message *sarama.ConsumerMessage, name string) string {
	nameBlob := convert.BytesFromString(name)

	for _, header := range message.Headers {
		if bytes.Equal(header.Key, nameBlob) {
			return convert.StringFromBytes(header.Value)
		}
	}

	return ""
}

// Headers returns all headers from message as map and [param.Param] object
func Headers(message *sarama.ConsumerMessage) map[string]string {
	headers := make(map[string]string, len(message.Headers))
	for _, header := range message.Headers {
		headers[string(header.Key)] = convert.StringFromBytes(header.Value)
	}
	return headers
}

// SetHeaders convert provided headers map to sarama headers slice
func SetHeaders(headers map[string]any) []sarama.RecordHeader {
	messageHeaders := make([]sarama.RecordHeader, len(headers))

	for name, value := range headers {
		messageHeaders = append(messageHeaders, sarama.RecordHeader{
			Key:   convert.BytesFromString(name),
			Value: convert.Bytes(value),
		})
	}

	return messageHeaders
}
