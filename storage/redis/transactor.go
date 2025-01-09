package redis

import (
	"context"
	"github.com/boostgo/lite/storage"
	"github.com/redis/go-redis/v9"
)

const (
	transactionKey = "lite_redis_tx"
)

type TransactorClientProvider interface {
	TxPipeline(ctx context.Context) (redis.Pipeliner, error)
}

type redisTransactor struct {
	provider TransactorClientProvider
}

func NewTransactor(provider TransactorClientProvider) storage.Transactor {
	return &redisTransactor{
		provider: provider,
	}
}

func (rt *redisTransactor) Begin(ctx context.Context) (storage.Transaction, error) {
	tx, err := rt.provider.TxPipeline(ctx)
	if err != nil {
		return nil, err
	}

	return newTransactorTx(ctx, tx), nil
}

func (rt *redisTransactor) BeginCtx(ctx context.Context) (context.Context, error) {
	tx, err := rt.provider.TxPipeline(ctx)
	if err != nil {
		return nil, err
	}

	return SetTx(ctx, tx), nil
}

func (rt *redisTransactor) CommitCtx(ctx context.Context) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil
	}

	_, err := tx.Exec(ctx)
	return err
}

func (rt *redisTransactor) RollbackCtx(ctx context.Context) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil
	}

	tx.Discard()
	return nil
}

func (rt *redisTransactor) Key() string {
	return transactionKey
}

type redisTransaction struct {
	tx        redis.Pipeliner
	parentCtx context.Context
}

func newTransactorTx(ctx context.Context, tx redis.Pipeliner) storage.Transaction {
	return &redisTransaction{
		tx:        tx,
		parentCtx: ctx,
	}
}

func (tx *redisTransaction) Commit(ctx context.Context) error {
	_, err := tx.tx.Exec(ctx)
	return err
}

func (tx *redisTransaction) Rollback(_ context.Context) error {
	tx.tx.Discard()
	return nil
}

func (tx *redisTransaction) Context() context.Context {
	return SetTx(tx.parentCtx, tx.tx)
}

// SetTx sets transaction key to new context
func SetTx(ctx context.Context, tx redis.Pipeliner) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

// GetTx returns sql transaction object from context if it exist
func GetTx(ctx context.Context) (redis.Pipeliner, bool) {
	transaction := ctx.Value(transactionKey)
	if transaction == nil {
		return nil, false
	}

	tx, ok := transaction.(redis.Pipeliner)
	return tx, ok
}
