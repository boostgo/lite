package kafka

import (
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func NewAsyncProducer(cfg Config, opts ...Option) (*AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	if cfg.Username != "" && cfg.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Handshake = true
		config.Net.SASL.Mechanism = "PLAIN"
		config.Net.SASL.User = cfg.Username
		config.Net.SASL.Password = cfg.Password
	}

	for _, opt := range opts {
		opt(config)
	}

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &AsyncProducer{
		producer: producer,
	}, nil
}

func (producer *AsyncProducer) Produce(messages ...*sarama.ProducerMessage) error {
	if len(messages) == 0 {
		return nil
	}

	for _, msg := range messages {
		producer.producer.Input() <- msg
	}

	return nil
}
