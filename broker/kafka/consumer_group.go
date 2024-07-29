package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"time"
)

type GroupHandler sarama.ConsumerGroupHandler
type GroupHandlerFunc func(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)

type ConsumerGroup struct {
	name   string
	group  sarama.ConsumerGroup
	topics []string
}

func NewConsumerGroup(name string, cfg Config, opts ...Option) (*ConsumerGroup, error) {
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
	life.Tear(consumerGroup.Close)

	return &ConsumerGroup{
		name:   name,
		group:  consumerGroup,
		topics: cfg.Topics,
	}, nil
}

func (consumer *ConsumerGroup) Close() error {
	return consumer.group.Close()
}

func (consumer *ConsumerGroup) Consume(handler GroupHandler) {
	logger := log.Namespace("kafka.consumer.group")

	go func() {
		for {
			select {
			case <-life.Context().Done():
				logger.Info().Str("name", consumer.name).Msg("Stop kafka consumer group")
				return
			default:
				if err := consumer.group.Consume(life.Context(), consumer.topics, handler); err != nil {
					logger.Error().Str("name", consumer.name).Err(err).Msg("Consume kafka handler")
					life.Cancel()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-consumer.group.Errors():
				logger.Error().Err(err).Str("name", consumer.name).Msg("Consumer group error")
				life.Cancel()
				return
			case <-life.Context().Done():
				logger.Info().Str("name", consumer.name).Msg("Stop worker from context")
				return
			}
		}
	}()
}

type ConsumerGroupHandlerFunc func(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
	message *sarama.ConsumerMessage,
)

type consumerGroupHandler struct {
	handler ConsumerGroupHandlerFunc
}

func ConsumerGroupHandler(handler ConsumerGroupHandlerFunc) sarama.ConsumerGroupHandler {
	return &consumerGroupHandler{
		handler: handler,
	}
}

func (handler *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return errors.New("kafka consumer channel closed")
			}

			handler.handler(session, claim, message)
		case <-session.Context().Done():
			return nil
		}
	}
}
