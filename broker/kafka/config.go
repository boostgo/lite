package kafka

import "github.com/IBM/sarama"

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
