package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
)

type SyncProducer struct {
	producer  sarama.SyncProducer
	traceMode bool
}

func SyncProducerOption() Option {
	return func(config *sarama.Config) {
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
	}
}

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

func MustSyncProducer(brokers []string, opts ...Option) *SyncProducer {
	producer, err := NewSyncProducer(brokers, opts...)
	if err != nil {
		panic(err)
	}

	return producer
}

func MustSyncProducerFromClient(client sarama.Client) *SyncProducer {
	producer, err := NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}

	return producer
}

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
