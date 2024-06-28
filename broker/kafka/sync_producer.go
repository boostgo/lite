package kafka

import (
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
)

type SyncProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(cfg *Config, opts ...Option) (*SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

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

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(producer.Close)

	return &SyncProducer{
		producer: producer,
	}, nil
}

func (producer *SyncProducer) Produce(messages ...*sarama.ProducerMessage) error {
	if len(messages) == 0 {
		return nil
	}

	return producer.producer.SendMessages(messages)
}
