package trace

import (
	"bytes"
	"context"

	"github.com/IBM/sarama"
	"github.com/boostgo/collection/slicex"
	"github.com/boostgo/lite/types/to"
)

// GetKafka return request id from kafka message
func GetKafka(message *sarama.ConsumerMessage) string {
	for _, header := range message.Headers {
		if to.String(header.Key) != _key {
			continue
		}

		return to.String(header.Value)
	}

	return ""
}

// GetKafkaCtx returns context with trace id from message
func GetKafkaCtx(ctx context.Context, message *sarama.ConsumerMessage) context.Context {
	for _, header := range message.Headers {
		if to.String(header.Key) != _key {
			continue
		}

		return Set(ctx, to.String(header.Value))
	}

	return ctx
}

// SetKafka sets trace id from context to provided messages
func SetKafka(ctx context.Context, messages ...*sarama.ProducerMessage) {
	if len(messages) == 0 {
		return
	}

	traceID := Get(ctx)
	if traceID == "" {
		return
	}

	traceIdBlob := to.Bytes(traceID)
	blobKey := to.Bytes(_key)
	for i := 0; i < len(messages); i++ {
		_, exist := slicex.Single(messages[i].Headers, func(header sarama.RecordHeader) bool {
			return bytes.Equal(header.Key, blobKey)
		})
		if exist {
			continue
		}

		messages[i].Headers = append(messages[i].Headers, sarama.RecordHeader{
			Key:   blobKey,
			Value: traceIdBlob,
		})
	}
}
