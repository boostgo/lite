package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/system/try"
)

type SnapshotClaim func(message *sarama.ConsumerMessage) error

func Snapshot(cfg Config, name, topic string, snapshotClaim SnapshotClaim, commit ...bool) error {
	offsets, err := GetOffsets(cfg.Brokers, sarama.NewConfig(), topic, sarama.OffsetNewest)
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

	consumer.consume(ctx, name, []string{topic}, ConsumerGroupHandler(
		name+" handler",
		func(ctx context.Context, session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, message *sarama.ConsumerMessage) error {
			lastOffset, ok := offsets[message.Partition]
			if !ok {
				return nil
			}

			if err = try.Try(func() error {
				return snapshotClaim(message)
			}); err != nil {
				return err
			}

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

			return nil
		},
		nil, nil,
	), func() {})

	<-done
	return nil
}
