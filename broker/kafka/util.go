package kafka

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/validator"
	"github.com/boostgo/lite/types/flex"
	"github.com/boostgo/lite/types/param"
	"github.com/boostgo/lite/types/to"
	"sync"
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

// GetOffsets return all offsets from topic
func GetOffsets(brokers []string, cfg *sarama.Config, topic string, offset int64) (map[int32]int64, error) {
	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	partitions, err := client.Partitions(topic)
	if err != nil {
		return nil, err
	}

	offsets := make(map[int32]int64)
	for i := 0; i < len(partitions); i++ {
		lastOffset, err := client.GetOffset(topic, partitions[i], offset)
		if err != nil {
			panic(err)
		}

		offsets[partitions[i]] = lastOffset - 1
	}

	return offsets, nil
}

// Parse message body to provided export object (which must be ptr) and validate for "validate" tags.
func Parse(message *sarama.ConsumerMessage, export any) error {
	// check export type
	if !flex.Type(export).IsPtr() {
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
func Header(message *sarama.ConsumerMessage, name string) param.Param {
	nameBlob := to.Bytes(name)

	for _, header := range message.Headers {
		if bytes.Equal(header.Key, nameBlob) {
			return param.New(to.String(header.Value))
		}
	}

	return param.Empty()
}

// Headers returns all headers from message as map and [param.Param] object
func Headers(message *sarama.ConsumerMessage) map[string]param.Param {
	headers := make(map[string]param.Param, len(message.Headers))
	for _, header := range message.Headers {
		headers[string(header.Key)] = param.New(to.String(header.Value))
	}
	return headers
}

// SetHeaders convert provided headers map to sarama headers slice
func SetHeaders(headers map[string]any) []sarama.RecordHeader {
	messageHeaders := make([]sarama.RecordHeader, len(headers))

	for name, value := range headers {
		messageHeaders = append(messageHeaders, sarama.RecordHeader{
			Key:   to.Bytes(name),
			Value: to.Bytes(value),
		})
	}

	return messageHeaders
}
