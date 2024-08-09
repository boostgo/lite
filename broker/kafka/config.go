package kafka

import (
	"errors"
	"github.com/IBM/sarama"
)

type Option func(*sarama.Config)

type Config struct {
	Brokers  []string
	Username string
	Password string
	GroupID  string
}

func With(fn func(*sarama.Config)) Option {
	return func(cfg *sarama.Config) {
		fn(cfg)
	}
}

func validateConsumerGroupConfig(config Config) error {
	if len(config.Brokers) == 0 {
		return errors.New("at least one broker is required")
	}

	if config.GroupID == "" {
		return errors.New("group id is required")
	}

	return nil
}

func validateConsumerConfig(config Config) error {
	if len(config.Brokers) == 0 {
		return errors.New("at least one broker is required")
	}

	return nil
}
