package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/storage"
	"github.com/jmoiron/sqlx"
)

func NotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

type DB interface {
	Connection() *sqlx.DB
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
	GetContext
	NamedExecContext
	SelectContext
	PrepareContext
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

type PrepareContext interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type client struct {
	conn      *sqlx.DB
	enableLog bool
}

func Client(conn *sqlx.DB, enableLog ...bool) DB {
	var enable bool
	if len(enableLog) > 0 {
		enable = enableLog[0]
	}

	return &client{
		conn:      conn,
		enableLog: enable,
	}
}

func (c *client) Connection() *sqlx.DB {
	return c.conn
}

func (c *client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	c.printLog(ctx, "ExecContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.conn.ExecContext(ctx, query, args...)
}

func (c *client) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	c.printLog(ctx, "QueryContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.conn.QueryContext(ctx, query, args...)
}

func (c *client) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	c.printLog(ctx, "QueryxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return c.conn.QueryxContext(ctx, query, args...)
}

func (c *client) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	c.printLog(ctx, "QueryRowxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return c.conn.QueryRowxContext(ctx, query, args...)
}

func (c *client) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	c.printLog(ctx, "PrepareContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareContext(ctx, query)
	}

	return c.conn.PrepareContext(ctx, query)
}

func (c *client) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	c.printLog(ctx, "NamedExecContext", query, arg)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return c.conn.NamedExecContext(ctx, query, arg)
}

func (c *client) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	c.printLog(ctx, "SelectContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return c.conn.SelectContext(ctx, dest, query, args...)
}

func (c *client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	c.printLog(ctx, "GetContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return c.conn.GetContext(ctx, dest, query, args...)
}

func (c *client) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	c.printLog(ctx, "PrepareNamedContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareNamedContext(ctx, query)
	}

	return c.conn.PrepareNamedContext(ctx, query)
}

func (c *client) printLog(ctx context.Context, queryType, query string, args ...any) {
	if !c.enableLog || storage.IsNoLog(ctx) {
		return
	}

	log.
		Context(ctx, "storage.sql."+queryType).
		Info().
		Str("query", query).
		Any("args", args).
		Send()
}

func Page(pageSize, page int) (offset int, limit int) {
	if page == 0 {
		page = 1
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, limit
}
