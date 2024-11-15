package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

func GetOffsets(brokers []string, cfg *sarama.Config, topic string, offset int64) (map[int32]int64, error) {
	client, err := sarama.NewClient(brokers, cfg)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	partitions, err := client.Partitions(topic)
	if err != nil {
		return nil, err
	}

	offsets := make(map[int32]int64)
	for i := 0; i < len(partitions); i++ {
		lastOffset, err := client.GetOffset(topic, partitions[i], offset)
		if err != nil {
			panic(err)
		}

		offsets[partitions[i]] = lastOffset - 1
	}

	return offsets, nil
}

func Parse(message *sarama.ConsumerMessage, export any) error {
	return json.Unmarshal(message.Value, export)
}
