package storage

import "context"

const (
	TransactionContextKey = "lite_tx"
)

// Transactor is common representation of transactions for any type of database.
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
