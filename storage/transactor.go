package storage

import (
	"context"
	"github.com/boostgo/lite/async"
)

const (
	TransactionContextKey = "lite_tx"
)

// Transactor is common representation of transactions for any type of database.
//
// Reason to use this: hide from usecase/service layer of using "sql" or "mongo" database
type Transactor interface {
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
func RewriteTx(original context.Context, toCopy context.Context) context.Context {
	tx := original.Value(TransactionContextKey)
	if tx == nil {
		return toCopy
	}

	return context.WithValue(toCopy, TransactionContextKey, tx)
}

type transactor struct {
	transactors []Transactor
}

func NewTransactor(transactors ...Transactor) Transactor {
	return &transactor{
		transactors: transactors,
	}
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
	tasks := make([]async.Task, 0, len(t.transactors))
	contexts := make(chan context.Context, len(t.transactors))
	for _, tr := range t.transactors {
		tasks = append(tasks, func() error {
			trCtx, err := tr.BeginCtx(ctx)
			if err != nil {
				return err
			}

			contexts <- trCtx
			return nil
		})
	}

	if err := async.WaitAll(tasks...); err != nil {
		return nil, err
	}

	close(contexts)

	ctxList := make([]context.Context, 0, len(t.transactors))
	for trCtx := range contexts {
		ctxList = append(ctxList, trCtx)
	}

	return newTransactorContext(ctxList...), nil
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
