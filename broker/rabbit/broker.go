package rabbit

import (
	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/system/trace"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Broker is wrap structure over amqp channel structure which collect all static functions in one object.
type Broker struct {
	channel       *amqp.Channel
	isTraceMaster bool
}

// NewBroker creates [Broker] with provided amqp channel
func NewBroker(channel *amqp.Channel) *Broker {
	return &Broker{
		channel:       channel,
		isTraceMaster: trace.AmIMaster(),
	}
}

// Close closes channel and connection
func (broker *Broker) Close() error {
	return broker.channel.Close()
}

// Ack acquires read message by delivery tag
func (broker *Broker) Ack(deliveryTag uint64, multiple bool) error {
	if err := broker.channel.Ack(deliveryTag, multiple); err != nil {
		return errorx.
			New("Ack message error").
			SetError(err).
			SetContext(map[string]any{
				"delivery-tag": deliveryTag,
				"multiple":     multiple,
			})
	}

	return nil
}

// Bind binds created exchange and queue.
func (broker *Broker) Bind(exchange, queue string) error {
	if err := broker.channel.QueueBind(
		queue,
		queue,
		exchange,
		false,
		nil,
	); err != nil {
		return errorx.
			New("Bind exchange and queue error").
			SetError(err).
			SetContext(map[string]any{
				"exchange": exchange,
				"queue":    queue,
			})
	}

	return nil
}
