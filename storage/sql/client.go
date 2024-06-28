package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

func NotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type DB interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
	GetContext
	NamedExecContext
	SelectContext
}

type NamedExecContext interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type SelectContext interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type GetContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type client struct {
	conn *sqlx.DB
}

func Client(conn *sqlx.DB) DB {
	return &client{
		conn: conn,
	}
}

func (c *client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.conn.ExecContext(ctx, query, args...)
}

func (c *client) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.conn.QueryContext(ctx, query, args...)
}

func (c *client) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return c.conn.QueryxContext(ctx, query, args...)
}

func (c *client) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return c.conn.QueryRowxContext(ctx, query, args...)
}

func (c *client) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareContext(ctx, query)
	}

	return c.conn.PrepareContext(ctx, query)
}

func (c *client) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return c.conn.NamedExecContext(ctx, query, arg)
}

func (c *client) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return c.conn.SelectContext(ctx, dest, query, args...)
}

func (c *client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, ok := GetTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return c.conn.GetContext(ctx, dest, query, args...)
}
