package kafka

import (
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/life"
	"github.com/google/uuid"
)

func NewClient(cfg Config, opts ...Option) (sarama.Client, error) {
	config := sarama.NewConfig()

	apply := make([]Option, 0, len(opts)+1)
	apply = append(apply, clientOption())
	apply = append(apply, opts...)

	for _, opt := range apply {
		opt(config)
	}

	client, err := sarama.NewClient(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}
	life.Tear(client.Close)

	return client, nil
}

func MustClient(cfg Config, opts ...Option) sarama.Client {
	client, err := NewClient(cfg, opts...)
	if err != nil {
		panic(err)
	}

	return client
}

func NewCluster(client sarama.Client) (sarama.ClusterAdmin, error) {
	clusterClient, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		return nil, err
	}

	return clusterClient, nil
}

func MustCluster(client sarama.Client) sarama.ClusterAdmin {
	cluster, err := NewCluster(client)
	if err != nil {
		panic(err)
	}

	return cluster
}

func clientOption() Option {
	return func(config *sarama.Config) {
		config.ClientID = buildClientID()
	}
}

func buildClientID() string {
	return "lite-app-" + uuid.New().String()
}
