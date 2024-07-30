package kafka

import (
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"time"
)

type ConsumeHandler func(message *sarama.ConsumerMessage) error
type ErrorHandler func(err error)

type Consumer struct {
	consumer     sarama.Consumer
	errorHandler ErrorHandler
}

func NewConsumer(cfg Config, opts ...Option) (*Consumer, error) {
	if err := validateConsumerConfig(cfg); err != nil {
		return nil, err
	}

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
		consumer: consumer,
	}, nil
}

func (consumer *Consumer) SetErrorHandler(handler ErrorHandler) {
	consumer.errorHandler = handler
}

func (consumer *Consumer) Consume(topic string, handler ConsumeHandler) error {
	logger := log.Namespace("kafka.consumer")

	partitions, err := consumer.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for i := 0; i < len(partitions); i++ {
		partitionConsumer, err := consumer.consumer.ConsumePartition(topic, partitions[i], sarama.OffsetNewest)
		if err != nil {
			return err
		}

		life.Tear(partitionConsumer.Close)

		go func(partition int32) {
			for {
				select {
				case <-life.Context().Done():
					logger.Info().Int32("partition", partition).Msg("Stop consumer by context")
					return
				case msg, ok := <-partitionConsumer.Messages():
					if !ok {
						logger.Info().Int32("partition", partition).Msg("Stop consumer by closing channel")
						return
					}

					if err = handler(msg); err != nil {
						logger.Error().Err(err).Msg("Handle message error")
						if consumer.errorHandler != nil {
							consumer.errorHandler(err)
						}
					}
				}
			}
		}(partitions[i])
	}

	return nil
}
