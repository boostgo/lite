package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
	"time"
)

type GroupHandler sarama.ConsumerGroupHandler
type GroupHandlerFunc func(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession)

type ConsumerGroup struct {
	group  sarama.ConsumerGroup
	topics []string
}

func NewConsumerGroup(cfg Config, opts ...Option) (*ConsumerGroup, error) {
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
		group:  consumerGroup,
		topics: cfg.Topics,
	}, nil
}

func (consumer *ConsumerGroup) Errors() <-chan error {
	return consumer.group.Errors()
}

func (consumer *ConsumerGroup) Close() error {
	return consumer.group.Close()
}

func (consumer *ConsumerGroup) Consume(handler GroupHandler) error {
	return consumer.group.Consume(life.Context(), consumer.topics, handler)
}

type groupHandler struct {
	ctx              context.Context
	groupHandlerFunc GroupHandlerFunc
}

func NewHandler(groupHandlerFunc GroupHandlerFunc) GroupHandler {
	return &groupHandler{
		ctx:              life.Context(),
		groupHandlerFunc: groupHandlerFunc,
	}
}

func (h *groupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-h.ctx.Done():
			return nil
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			h.groupHandlerFunc(message, session)
		}
	}
}
