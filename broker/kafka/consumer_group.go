package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/collections/list"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"time"
)

type GroupHandler sarama.ConsumerGroupHandler
type GroupHandlerFunc func(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)

// ConsumerGroup wrap structure for [Consumer] Group
type ConsumerGroup struct {
	group sarama.ConsumerGroup
}

// ConsumerGroupOption returns default consumer group configs
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

// NewConsumerGroup creates [ConsumerGroup] by options
func NewConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
	consumerGroup, err := newConsumerGroup(cfg, opts...)
	if err != nil {
		return nil, err
	}
	life.Tear(consumerGroup.Close)

	return consumerGroup, nil
}

// NewConsumerGroupFromClient creates [ConsumerGroup] by sarama client
func NewConsumerGroupFromClient(groupID string, client sarama.Client) (*ConsumerGroup, error) {
	consumerGroup, err := newConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, err
	}
	life.Tear(consumerGroup.Close)

	return consumerGroup, nil
}

// MustConsumerGroup calls [NewConsumerGroup] and if error catch throws panic
func MustConsumerGroup(cfg Config, opts ...Option) *ConsumerGroup {
	consumer, err := NewConsumerGroup(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return consumer
}

// MustConsumerGroupFromClient calls [NewConsumerGroupFromClient] and if error catch throws panic
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

// Consume starts consuming topic with consumer group.
// Catch consumer group errors and provided context done (for graceful shutdown).
func (consumer *ConsumerGroup) Consume(name string, topics []string, handler GroupHandler) {
	consumer.consume(life.Context(), name, topics, handler, life.Cancel)
}

func (consumer *ConsumerGroup) consume(
	ctx context.Context,
	name string,
	topics []string,
	handler GroupHandler,
	cancel func(),
) {
	logger := log.Namespace("kafka.consumer.group")

	// run consuming
	go func() {
		if err := consumer.group.Consume(ctx, topics, handler); err != nil {
			logger.
				Error().
				Err(err).
				Str("name", name).
				Strs("topics", topics).
				Msg("Consume kafka claim")
			cancel()
		}
	}()

	// run catching consumer errors and context canceling
	go func() {
		for {
			select {
			case err := <-consumer.group.Errors():
				logger.
					Error().
					Err(err).
					Str("name", name).
					Strs("topics", topics).
					Msg("Consumer group error")
				cancel()
				return
			case <-ctx.Done():
				logger.
					Info().
					Str("name", name).
					Strs("topics", topics).
					Msg("Stop broker from context")
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
	name    string
	claim   ConsumerGroupClaim
	setup   ConsumerGroupSetup
	cleanup ConsumerGroupCleanup
	timeout time.Duration
}

// ConsumerGroupHandler creates [sarama.ConsumerGroupHandler] interface implementation object.
// Provide 3 methods: claim, setup and cleanup for implementing interface.
// Also, could be provided timeout for claiming every message
func ConsumerGroupHandler(
	name string,
	handler ConsumerGroupClaim,
	setup ConsumerGroupSetup,
	cleanup ConsumerGroupCleanup,
	timeout ...time.Duration,
) sarama.ConsumerGroupHandler {
	var setTimeout time.Duration
	if len(timeout) > 0 {
		setTimeout = timeout[0]
	}

	return &consumerGroupHandler{
		name:    name,
		claim:   handler,
		setup:   setup,
		cleanup: cleanup,
		timeout: setTimeout,
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
				return errs.
					New("Kafka Consumer Group Handler channel closed").
					AddContext("name", handler.name)
			}

			func() {
				ctx := context.Background()
				var cancel context.CancelFunc

				if handler.timeout > 0 {
					ctx, cancel = context.WithTimeout(ctx, handler.timeout)
					defer cancel()
				}

				ctx = trace.GetKafkaCtx(ctx, message)

				if err := try.Try(func() error {
					return handler.claim(ctx, session, claim, message)
				}); err != nil {
					logger.
						Error().
						Ctx(ctx).
						Err(err).
						Str("name", handler.name).
						Msg("Kafka consumer group claim")
				}
			}()
		case <-session.Context().Done():
			return nil
		}
	}
}
