package trace

import (
	"context"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// GetAMQP returns trace id from RMQ message
func GetAMQP(message *amqp.Delivery) string {
	if message == nil || message.Headers == nil {
		return ""
	}

	traceID := message.Headers[_key]
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

// SetAmqp set trace id in RMQ headers from context
func SetAmqp(table amqp.Table, ctx context.Context) {
	traceID := Get(ctx)
	if traceID == "" {
		return
	}

	table[_key] = traceID
}

// SetAmqpCtx set trace id to RMQ headers table
func SetAmqpCtx(ctx context.Context, table amqp.Table) bool {
	traceID := Get(ctx)
	if traceID == "" {
		return false
	}

	table[_key] = traceID
	return true
}
