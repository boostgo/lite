package kafka

import (
	"errors"
	"github.com/IBM/sarama"
)

type ConsumerGroupHandlerFunc func(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
	message *sarama.ConsumerMessage,
)

type ConsumerGroupHandler struct {
	handler ConsumerGroupHandlerFunc
}

func NewConsumerGroupHandler(handler ConsumerGroupHandlerFunc) sarama.ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		handler: handler,
	}
}

func (handler *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (handler *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return errors.New("kafka consumer channel closed")
			}

			handler.handler(session, claim, message)
		case <-session.Context().Done():
			return nil
		}
	}
}
