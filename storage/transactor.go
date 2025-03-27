package storage

import (
	"context"

	"github.com/boostgo/collection/slicex"
	"github.com/boostgo/lite/async"
)

// Transactor is common representation of transactions for any type of database.
//
// Reason to use this: hide from usecase/service layer of using "sql" or "mongo" database
type Transactor interface {
	Key() string
	IsTx(ctx context.Context) bool
	Begin(ctx context.Context) (Transaction, error)
	BeginCtx(ctx context.Context) (context.Context, error)
	CommitCtx(ctx context.Context) error
	RollbackCtx(ctx context.Context) error
}

// Transaction interface using by Transactor
type Transaction interface {
	Context() context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// RewriteTx take transaction key from original context and copy key to toCopy context
func RewriteTx(key string, original context.Context, toCopy context.Context) context.Context {
	tx := original.Value(key)
	if tx == nil {
		return toCopy
	}

	return context.WithValue(toCopy, key, tx)
}

type transactor struct {
	transactors []Transactor
}

func NewTransactor(transactors ...Transactor) Transactor {
	return &transactor{
		transactors: transactors,
	}
}

func (t *transactor) Key() string {
	return slicex.
		JoinString(slicex.Map(t.transactors, func(t Transactor) string {
			return t.Key()
		}), func(s string) string {
			return s
		})
}

func (t *transactor) IsTx(ctx context.Context) bool {
	for _, tr := range t.transactors {
		if tr.IsTx(ctx) {
			return true
		}
	}

	return false
}

func (t *transactor) Begin(ctx context.Context) (Transaction, error) {
	transactions := make([]Transaction, 0, len(t.transactors))
	for _, tr := range t.transactors {
		tx, err := tr.Begin(ctx)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}

	return newTransaction(transactions), nil
}

func (t *transactor) BeginCtx(ctx context.Context) (context.Context, error) {
	var err error
	for _, tr := range t.transactors {
		ctx, err = tr.BeginCtx(ctx)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

func (t *transactor) CommitCtx(ctx context.Context) error {
	tasks := make([]async.Task, 0, len(t.transactors))
	for _, tr := range t.transactors {
		tasks = append(tasks, func() error {
			return tr.CommitCtx(ctx)
		})
	}

	return async.WaitAll(tasks...)
}

func (t *transactor) RollbackCtx(ctx context.Context) error {
	tasks := make([]async.Task, 0, len(t.transactors))
	for _, tr := range t.transactors {
		tasks = append(tasks, func() error {
			return tr.RollbackCtx(ctx)
		})
	}

	return async.WaitAll(tasks...)
}

type transaction struct {
	transactions []Transaction
}

func (t *transaction) Context() context.Context {
	return nil
}

func (t *transaction) Commit(ctx context.Context) error {
	tasks := make([]async.Task, 0, len(t.transactions))
	for _, tx := range t.transactions {
		tasks = append(tasks, func() error {
			return tx.Commit(ctx)
		})
	}

	return async.WaitAll(tasks...)
}

func (t *transaction) Rollback(ctx context.Context) error {
	tasks := make([]async.Task, 0, len(t.transactions))
	for _, tx := range t.transactions {
		tasks = append(tasks, func() error {
			return tx.Rollback(ctx)
		})
	}

	return async.WaitAll(tasks...)
}

func newTransaction(transactions []Transaction) Transaction {
	return &transaction{
		transactions: transactions,
	}
}
