package kafka

import (
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func AsyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}
}

func NewAsyncProducer(brokers []string, opts ...Option) (*AsyncProducer, error) {
	config := sarama.NewConfig()
	config.ClientID = buildClientID()
	AsyncProducerOption()(config)

	for _, opt := range opts {
		opt(config)
	}

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &AsyncProducer{
		producer: producer,
	}, nil
}

func NewAsyncProducerFromClient(client sarama.Client) (*AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &AsyncProducer{
		producer: producer,
	}, nil
}

func MustAsyncProducer(brokers []string, opts ...Option) *AsyncProducer {
	producer, err := NewAsyncProducer(brokers, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

func MustAsyncProducerFromClient(client sarama.Client) *AsyncProducer {
	producer, err := NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
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
