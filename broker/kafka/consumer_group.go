package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/try"
	"time"
)

type GroupHandler sarama.ConsumerGroupHandler
type GroupHandlerFunc func(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)

type ConsumerGroup struct {
	group  sarama.ConsumerGroup
	topics []string
}

func NewConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
	consumerGroup, err := newConsumerGroup(cfg, opts...)
	if err != nil {
		return nil, err
	}
	life.Tear(consumerGroup.Close)

	return consumerGroup, nil
}

func newConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
	if err := validateConsumerGroupConfig(cfg); err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second

	config.Consumer.Group.Rebalance.GroupStrategies = append(
		config.Consumer.Group.Rebalance.GroupStrategies, sarama.NewBalanceStrategyRoundRobin())
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

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group:  consumerGroup,
		topics: cfg.Topics,
	}, nil
}

func (consumer *ConsumerGroup) Close() error {
	return consumer.group.Close()
}

func (consumer *ConsumerGroup) Consume(name string, handler GroupHandler) {
	consumer.consume(life.Context(), name, handler, life.Cancel)
}

func (consumer *ConsumerGroup) consume(ctx context.Context, name string, handler GroupHandler, cancel func()) {
	logger := log.Namespace("kafka.consumer.group")

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info().Str("name", name).Msg("Stop kafka consumer group")
				return
			default:
				if err := consumer.group.Consume(life.Context(), consumer.topics, handler); err != nil {
					logger.Error().Str("name", name).Err(err).Msg("Consume kafka claim")
					cancel()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-consumer.group.Errors():
				logger.Error().Err(err).Str("name", name).Msg("Consumer group error")
				cancel()
				return
			case <-ctx.Done():
				logger.Info().Str("name", name).Msg("Stop worker from context")
				return
			}
		}
	}()
}

type (
	ConsumerGroupClaim func(
		session sarama.ConsumerGroupSession,
		claim sarama.ConsumerGroupClaim,
		message *sarama.ConsumerMessage,
	) error

	ConsumerGroupSetup   func(session sarama.ConsumerGroupSession) error
	ConsumerGroupCleanup func(session sarama.ConsumerGroupSession) error
)

type consumerGroupHandler struct {
	claim   ConsumerGroupClaim
	setup   ConsumerGroupSetup
	cleanup ConsumerGroupCleanup
}

func ConsumerGroupHandler(
	handler ConsumerGroupClaim,
	setup ConsumerGroupSetup,
	cleanup ConsumerGroupCleanup,
) sarama.ConsumerGroupHandler {
	return &consumerGroupHandler{
		claim:   handler,
		setup:   setup,
		cleanup: cleanup,
	}
}

func (handler *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	if handler.setup != nil {
		return handler.setup(session)
	}

	return nil
}

func (handler *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	if handler.cleanup != nil {
		return handler.cleanup(session)
	}

	return nil
}

func (handler *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger := log.Namespace("kafka.consumer.group")

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return errors.New("kafka consumer channel closed")
			}

			if err := try.Try(func() error {
				return handler.claim(session, claim, message)
			}); err != nil {
				logger.Error().Err(err).Msg("Kafka consumer group claim")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
