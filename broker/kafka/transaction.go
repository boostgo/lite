package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/boostgo/lite/storage"
	"github.com/boostgo/lite/types/to"
)

var _ storage.Transaction = new(kafkaTransaction)
var _ storage.Transactor = new(kafkaTransactor)

const (
	txKey   = "kafka_tx"
	noTxKey = "kafka_no_tx"
)

type kafkaTransactor struct {
	producer *SyncProducer
}

func NewTransactor(producer *SyncProducer) storage.Transactor {
	return &kafkaTransactor{
		producer: producer,
	}
}

func (t *kafkaTransactor) Key() string {
	return txKey
}

func (t *kafkaTransactor) IsTx(ctx context.Context) bool {
	txObject := ctx.Value(txKey)
	if txObject == nil {
		return false
	}

	_, ok := txObject.(*kafkaTransaction)
	if !ok {
		return false
	}

	return true
}

func (t *kafkaTransactor) Begin(_ context.Context) (storage.Transaction, error) {
	return newTransaction(t.producer), nil
}

func (t *kafkaTransactor) BeginCtx(ctx context.Context) (context.Context, error) {
	tx := newTransaction(t.producer)
	return context.WithValue(ctx, txKey, tx), nil
}

func (t *kafkaTransactor) CommitCtx(ctx context.Context) error {
	tx, ok := getTx(ctx)
	if !ok {
		return nil
	}

	return tx.Commit(ctx)
}

func (t *kafkaTransactor) RollbackCtx(ctx context.Context) error {
	tx, ok := getTx(ctx)
	if !ok {
		return nil
	}

	return tx.Rollback(ctx)
}

type kafkaTransaction struct {
	producer *SyncProducer
	messages []*sarama.ProducerMessage
}

func newTransaction(producer *SyncProducer) storage.Transaction {
	return &kafkaTransaction{
		producer: producer,
		messages: make([]*sarama.ProducerMessage, 0),
	}
}

func (tx *kafkaTransaction) addMessages(messages ...*sarama.ProducerMessage) {
	if len(messages) == 0 {
		return
	}

	tx.messages = append(tx.messages, messages...)
}

func (tx *kafkaTransaction) Context() context.Context {
	return context.Background()
}

func (tx *kafkaTransaction) Commit(ctx context.Context) error {
	if tx.messages == nil {
		return nil
	}

	return tx.producer.Produce(context.Background(), tx.messages...)
}

func (tx *kafkaTransaction) Rollback(_ context.Context) error {
	tx.messages = nil
	return nil
}

func NoTxContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, noTxKey, true)
}

func getTx(ctx context.Context) (storage.Transaction, bool) {
	noTx := ctx.Value(noTxKey)
	if noTx != nil && to.Bool(noTx) {
		return nil, false
	}

	txObject := ctx.Value(txKey)
	if txObject == nil {
		return nil, false
	}

	tx, ok := txObject.(storage.Transaction)
	if !ok {
		return nil, false
	}

	return tx, true
}

func addMessagesToTx(tx storage.Transaction, messages []*sarama.ProducerMessage) bool {
	kafkaTx, ok := tx.(*kafkaTransaction)
	if !ok {
		return false
	}

	kafkaTx.addMessages(messages...)
	return true
}
