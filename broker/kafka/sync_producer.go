package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
)

// SyncProducer producer which produce messages in current goroutine
type SyncProducer struct {
	producer  sarama.SyncProducer
	traceMode bool
}

// SyncProducerOption returns default sync producer configuration as [Option]
func SyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}
}

// NewSyncProducer creates [SyncProducer] with configurations.
//
// Creates sync producer with default configuration as [Option] created by [SyncProducerOption] function.
//
// Adds producer close method to teardown
func NewSyncProducer(brokers []string, opts ...Option) (*SyncProducer, error) {
	config := sarama.NewConfig()
	config.ClientID = buildClientID()
	SyncProducerOption()(config)

	for _, opt := range opts {
		opt(config)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &SyncProducer{
		producer:  producer,
		traceMode: trace.AmIMaster(),
	}, nil
}

// NewSyncProducerFromClient creates [SyncProducer] by provided client.
//
// Creates sync producer with default configuration as Option created by SyncProducerOption function.
//
// Adds producer close method to teardown
func NewSyncProducerFromClient(client sarama.Client) (*SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &SyncProducer{
		producer:  producer,
		traceMode: trace.AmIMaster(),
	}, nil
}

// MustSyncProducer calls [NewSyncProducer] function with calls panic if returns error
func MustSyncProducer(brokers []string, opts ...Option) *SyncProducer {
	producer, err := NewSyncProducer(brokers, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

// MustSyncProducerFromClient calls [NewSyncProducerFromClient] function with calls panic if returns error
func MustSyncProducerFromClient(client sarama.Client) *SyncProducer {
	producer, err := NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

// Produce sends provided message(s) in the same goroutine.
//
// Sets trace id to provided messages to header
func (producer *SyncProducer) Produce(ctx context.Context, messages ...*sarama.ProducerMessage) error {
	if len(messages) == 0 {
		return nil
	}

	if producer.traceMode && trace.Get(ctx) == "" {
		ctx = trace.Set(ctx, trace.String())
	}

	trace.SetKafka(ctx, messages...)
	return producer.producer.SendMessages(messages)
}
