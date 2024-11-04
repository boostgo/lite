package sql

import (
	"context"
	"database/sql"
	"github.com/boostgo/lite/storage"
	"github.com/jmoiron/sqlx"
)

type TransactorConnectionProvider interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type sqlTransactor struct {
	provider TransactorConnectionProvider
}

func NewTransactor(provider TransactorConnectionProvider) storage.Transactor {
	return &sqlTransactor{
		provider: provider,
	}
}

func (st sqlTransactor) Begin(ctx context.Context) (storage.Transaction, error) {
	tx, err := st.provider.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	return newTransactorTx(ctx, tx), nil
}

func (st sqlTransactor) BeginCtx(ctx context.Context) (context.Context, error) {
	tx, err := st.provider.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	return SetTx(ctx, tx), nil
}

func (st sqlTransactor) CommitCtx(ctx context.Context) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil
	}

	return tx.Commit()
}

func (st sqlTransactor) RollbackCtx(ctx context.Context) error {
	tx, ok := GetTx(ctx)
	if !ok {
		return nil
	}

	return tx.Rollback()
}

type sqlTransaction struct {
	tx        *sqlx.Tx
	parentCtx context.Context
}

func newTransactorTx(ctx context.Context, tx *sqlx.Tx) storage.Transaction {
	return &sqlTransaction{
		tx:        tx,
		parentCtx: ctx,
	}
}

func (tx sqlTransaction) Commit(_ context.Context) error {
	return tx.tx.Commit()
}

func (tx sqlTransaction) Rollback(_ context.Context) error {
	return tx.tx.Rollback()
}

func (tx sqlTransaction) Context() context.Context {
	return SetTx(tx.parentCtx, tx.tx)
}
