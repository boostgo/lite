package sql

import (
	"context"
	"github.com/boostgo/lite/system/try"
	"github.com/jmoiron/sqlx"
)

const (
	txKey = "lite_tx"
)

func SetTx(conn *sqlx.DB, ctx context.Context) (context.Context, error) {
	transaction, err := conn.Beginx()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, txKey, transaction), nil
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	transaction := ctx.Value(txKey)
	if transaction == nil {
		return nil, false
	}

	tx, ok := transaction.(*sqlx.Tx)
	return tx, ok
}

func Transaction(conn *sqlx.DB, transactionActions func(tx *sqlx.Tx) error) error {
	transaction, err := conn.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = transaction.Rollback()
	}()

	if err = transactionActions(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}

func Atomic(ctx context.Context, conn *sqlx.DB, fn func(ctx context.Context) error) error {
	tx, err := conn.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	defer func() {
		_ = tx.Commit()
	}()

	return try.Try(func() error {
		return fn(context.WithValue(ctx, txKey, tx))
	})
}
