package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/collections/list"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"time"
)

type GroupHandler sarama.ConsumerGroupHandler
type GroupHandlerFunc func(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)

type ConsumerGroup struct {
	group sarama.ConsumerGroup
}

func ConsumerGroupOption(offset ...int64) Option {
	return func(config *sarama.Config) {
		config.Consumer.Return.Errors = true

		if len(offset) > 0 {
			config.Consumer.Offsets.Initial = offset[0]
		} else {
			config.Consumer.Offsets.Initial = sarama.OffsetNewest
		}

		config.Consumer.Offsets.AutoCommit.Enable = true
		config.Consumer.Offsets.AutoCommit.Interval = time.Second

		config.Consumer.Group.Rebalance.GroupStrategies = append(
			config.Consumer.Group.Rebalance.GroupStrategies, sarama.NewBalanceStrategyRoundRobin())
		config.Consumer.Fetch.Default = 1 << 20 // 1MB
		config.Consumer.Fetch.Max = 10 << 20    // 10MB
		config.ChannelBufferSize = 256
	}
}

func NewConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
	consumerGroup, err := newConsumerGroup(cfg, opts...)
	if err != nil {
		return nil, err
	}
	life.Tear(consumerGroup.Close)

	return consumerGroup, nil
}

func NewConsumerGroupFromClient(groupID string, client sarama.Client) (*ConsumerGroup, error) {
	consumerGroup, err := newConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}
	life.Tear(consumerGroup.Close)

	return consumerGroup, nil
}

func MustConsumerGroup(cfg Config, opts ...Option) *ConsumerGroup {
	consumer, err := NewConsumerGroup(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return consumer
}

func MustConsumerGroupFromClient(groupID string, client sarama.Client) *ConsumerGroup {
	consumer, err := NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err)
	}

	return consumer
}

func newConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
	client, err := NewClient(cfg, list.AddLeft(opts, ConsumerGroupOption())...)
	if err != nil {
		return nil, err
	}

	return newConsumerGroupFromClient(cfg.GroupID, client)
}

func newConsumerGroupFromClient(groupID string, client sarama.Client) (*ConsumerGroup, error) {
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		group: consumerGroup,
	}, nil
}

func (consumer *ConsumerGroup) Close() error {
	return consumer.group.Close()
}

func (consumer *ConsumerGroup) Consume(name string, topics []string, handler GroupHandler) {
	consumer.consume(life.Context(), name, topics, handler, life.Cancel)
}

func (consumer *ConsumerGroup) consume(ctx context.Context, name string, topics []string, handler GroupHandler, cancel func()) {
	logger := log.Namespace("kafka.consumer.group")

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info().Str("name", name).Msg("Stop kafka consumer group")
				return
			default:
				if err := consumer.group.Consume(ctx, topics, handler); err != nil {
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
				logger.Info().Str("name", name).Msg("Stop broker from context")
				return
			}
		}
	}()
}

type (
	ConsumerGroupClaim func(
		ctx context.Context,
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
	logger := log.Namespace("kafka.consumer.group.handler")

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return errors.New("kafka consumer channel closed")
			}

			func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
				defer cancel()

				ctx = trace.GetKafkaCtx(ctx, message)

				if err := try.Try(func() error {
					return handler.claim(ctx, session, claim, message)
				}); err != nil {
					logger.Error().Err(err).Msg("Kafka consumer group claim")
				}
			}()
		case <-session.Context().Done():
			return nil
		}
	}
}
