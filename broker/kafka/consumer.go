package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"time"
)

type ConsumeHandler func(message *sarama.ConsumerMessage) error
type ErrorHandler func(err error)

type Consumer struct {
	l            log.Logger
	ctx          context.Context
	consumer     sarama.Consumer
	topic        string
	handler      ConsumeHandler
	errorHandler ErrorHandler
}

func NewConsumer(handler ConsumeHandler, cfg Config, opts ...Option) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second

	config.Consumer.Fetch.Default = 1 << 20 // 1MB
	config.Consumer.Fetch.Max = 10 << 20    // 10MB
	config.ChannelBufferSize = 256

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

	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(consumer.Close)

	return &Consumer{
		l:        log.Namespace("kafka.consumer"),
		ctx:      life.Context(),
		consumer: consumer,
		handler:  handler,
	}, nil
}

func (consumer *Consumer) SetErrorHandler(handler ErrorHandler) {
	consumer.errorHandler = handler
}

func (consumer *Consumer) Consume() error {
	partitions, err := consumer.consumer.Partitions(consumer.topic)
	if err != nil {
		return err
	}

	for i := 0; i < len(partitions); i++ {
		partitionConsumer, err := consumer.consumer.ConsumePartition(consumer.topic, partitions[i], sarama.OffsetNewest)
		if err != nil {
			return err
		}

		life.Tear(partitionConsumer.Close)

		go func() {
			select {
			case <-consumer.ctx.Done():
				consumer.l.Info().Int("partition", i).Msg("Stop consumer by context")
				return
			case msg, ok := <-partitionConsumer.Messages():
				if !ok {
					consumer.l.Info().Int("partition", i).Msg("Stop consumer by closing channel")
					return
				}

				if err = consumer.handler(msg); err != nil {
					consumer.l.Error().Err(err).Msg("Handle message error")
					if consumer.errorHandler != nil {
						consumer.errorHandler(err)
					}
				}
			}
		}()
	}

	return nil
}
