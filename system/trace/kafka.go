package trace

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/types/to"
)

func GetKafka(message *sarama.ConsumerMessage) string {
	for _, header := range message.Headers {
		if to.String(header.Key) != key {
			continue
		}

		return to.String(header.Value)
	}

	return ""
}

func GetKafkaCtx(ctx context.Context, message *sarama.ConsumerMessage) context.Context {
	for _, header := range message.Headers {
		if to.String(header.Key) != key {
			continue
		}

		return Set(ctx, to.String(header.Value))
	}

	return ctx
}

func SetKafka(ctx context.Context, messages ...*sarama.ProducerMessage) {
	if len(messages) == 0 {
		return
	}

	traceID := Get(ctx)
	if traceID == "" {
		return
	}

	traceIdBlob := to.Bytes(traceID)
	for i := 0; i < len(messages); i++ {
		messages[i].Headers = append(messages[i].Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: traceIdBlob,
		})
	}
}