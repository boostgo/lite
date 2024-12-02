package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
)

// AsyncProducer producer which produce messages in "async" way
type AsyncProducer struct {
	producer  sarama.AsyncProducer
	traceMode bool
}

// AsyncProducerOption returns default async producer configuration as [Option]
func AsyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}
}

// NewAsyncProducer creates [AsyncProducer] with configurations.
// Creates async producer with default configuration as [Option] created by [AsyncProducerOption] function.
// Adds producer close method to teardown
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
		producer:  producer,
		traceMode: trace.AmIMaster(),
	}, nil
}

// NewAsyncProducerFromClient creates [AsyncProducer] by provided client.
// Creates async producer with default configuration as [Option] created by [AsyncProducerOption] function.
// Adds producer close method to teardown
func NewAsyncProducerFromClient(client sarama.Client) (*AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &AsyncProducer{
		producer:  producer,
		traceMode: trace.AmIMaster(),
	}, nil
}

// MustAsyncProducer calls [NewAsyncProducer] function with calls panic if returns error
func MustAsyncProducer(brokers []string, opts ...Option) *AsyncProducer {
	producer, err := NewAsyncProducer(brokers, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

// MustAsyncProducerFromClient calls [NewAsyncProducerFromClient] function with calls panic if returns error
func MustAsyncProducerFromClient(client sarama.Client) *AsyncProducer {
	producer, err := NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

// Produce sends provided message(s) in other goroutine.
// Sets trace id to provided messages to header
func (producer *AsyncProducer) Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error {
	if len(messages) == 0 {
		return nil
	}

	if producer.traceMode && trace.Get(ctx) == "" {
		ctx = trace.Set(ctx, trace.String())
	}

	trace.SetKafka(ctx, messages...)

	for _, msg := range messages {
		producer.producer.Input() <- msg
	}

	return nil
}
