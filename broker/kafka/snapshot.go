package kafka

import (
	"context"
	"github.com/IBM/sarama"
)

type SnapshotClaim func(message *sarama.ConsumerMessage)

func Snapshot(cfg Config, name string, brokers []string, topic string, snapshotClaim SnapshotClaim, commit ...bool) error {
	offsets, err := GetOffsets(brokers, sarama.NewConfig(), topic, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	consumer, err := newConsumerGroup(cfg, func(config *sarama.Config) {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	})
	if err != nil {
		return err
	}

	done := make(chan struct{}, 1)

	ctx, cancel := context.WithCancel(context.Background())

	var autoCommit bool
	if len(commit) > 0 {
		autoCommit = commit[0]
	}

	consumer.consume(ctx, name, ConsumerGroupHandler(
		func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, message *sarama.ConsumerMessage) {
			lastOffset, ok := offsets[message.Partition]
			if !ok {
				return
			}

			snapshotClaim(message)
			if autoCommit {
				session.MarkMessage(message, "")
			}

			if lastOffset == message.Offset {
				delete(offsets, message.Partition)
				if len(offsets) == 0 {
					cancel()
					go consumer.group.Close()
					done <- struct{}{}
				}
			}
		},
		nil, nil,
	), func() {})

	<-done
	return nil
}
