package storage

import "context"

// Transactor is common representation of transactions for any type of database.
// Reason to use this: hide from usecase/service layer of using "sql" or "mongo" database
type Transactor interface {
	Begin(ctx context.Context) (Transaction, error)
	BeginCtx(ctx context.Context) (context.Context, error)
	CommitCtx(ctx context.Context) error
	RollbackCtx(ctx context.Context) error
}

type Transaction interface {
	Context() context.Context
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
