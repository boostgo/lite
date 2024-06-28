package rabbit

import (
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/trace"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	channel       *amqp.Channel
	isTraceMaster bool
}

func NewBroker(channel *amqp.Channel) *Broker {
	return &Broker{
		channel:       channel,
		isTraceMaster: trace.AmIMaster(),
	}
}

func (broker *Broker) Close() error {
	return broker.channel.Close()
}

func (broker *Broker) Ack(deliveryTag uint64, multiple bool) error {
	if err := broker.channel.Ack(deliveryTag, multiple); err != nil {
		return errs.
			New("Ack message error").
			SetError(err).
			SetContext(map[string]any{
				"delivery-tag": deliveryTag,
				"multiple":     multiple,
			})
	}

	return nil
}

func (broker *Broker) Bind(exchange, queue string) error {
	if err := broker.channel.QueueBind(
		queue,
		queue,
		exchange,
		false,
		nil,
	); err != nil {
		return errs.
			New("Bind exchange and queue error").
			SetError(err).
			SetContext(map[string]any{
				"exchange": exchange,
				"queue":    queue,
			})
	}

	return nil
}
