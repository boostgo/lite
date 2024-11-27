package rabbit

import (
	"github.com/boostgo/lite/errs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagesQueue <-chan amqp.Delivery

type ConsumerConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func defaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{}
}

// Consume starts consuming provided queue with optional configurations
func (broker *Broker) Consume(queue string, cfg ...ConsumerConfig) (<-chan amqp.Delivery, error) {
	var config ConsumerConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = defaultConsumerConfig()
	}

	messages, err := broker.channel.Consume(
		queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return nil, errs.
			New("Consume queue error").
			SetError(err).
			AddContext("queue", queue)
	}

	return messages, nil
}
