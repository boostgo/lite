package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"strconv"
)

type Option func(*sarama.Config)

type Config struct {
	Brokers  []string
	Username string
	Password string

	Topics  []string
	GroupID string
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

	if len(config.Topics) == 0 {
		return errors.New("at least one topic is required")
	}

	if config.GroupID == "" {
		return errors.New("group id is required")
	}

	for i := 0; i < len(config.Topics); i++ {
		if config.Topics[i] == "" {
			return errors.New("one of topics is empty: " + strconv.Itoa(i))
		}
	}

	return nil
}

func validateConsumerConfig(config Config) error {
	if len(config.Brokers) == 0 {
		return errors.New("at least one broker is required")
	}

	if len(config.Topics) == 0 {
		return errors.New("at least one topic is required")
	}

	for i := 0; i < len(config.Topics); i++ {
		if config.Topics[i] == "" {
			return errors.New("one of topics is empty: " + strconv.Itoa(i))
		}
	}

	return nil
}
