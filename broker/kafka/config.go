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

func (cfg Config) Copy(groupID string) Config {
	return Config{
		Brokers:  cfg.Brokers,
		Username: cfg.Username,
		Password: cfg.Password,
		GroupID:  groupID,
	}
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
