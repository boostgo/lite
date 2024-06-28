package trace

import (
	"context"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetAMQP(message *amqp.Delivery) string {
	if message == nil || message.Headers == nil {
		return ""
	}

	traceID := message.Headers[key]
	if traceID == nil {
		return ""
	}

	switch tid := traceID.(type) {
	case string:
		return tid
	case uuid.UUID:
		return tid.String()
	default:
		return ""
	}
}

func SetAmqp(table amqp.Table, traceID string) {
	table[key] = traceID
}

func SetAmqpCtx(ctx context.Context, table amqp.Table) bool {
	traceID := Get(ctx)
	if traceID == "" {
		return false
	}

	table[key] = traceID
	return true
}
