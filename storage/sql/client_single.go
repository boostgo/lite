package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/storage"
	"github.com/jmoiron/sqlx"
)

type clientSingle struct {
	conn      *sqlx.DB
	enableLog bool
}

func Client(conn *sqlx.DB, enableLog ...bool) DB {
	var enable bool
	if len(enableLog) > 0 {
		enable = enableLog[0]
	}

	return &clientSingle{
		conn:      conn,
		enableLog: enable,
	}
}

func (c *clientSingle) Connection() *sqlx.DB {
	return c.conn
}

func (c *clientSingle) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	defer errs.Wrap(errType, &err, "ExecContext")

	c.printLog(ctx, "ExecContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.conn.ExecContext(ctx, query, args...)
}

func (c *clientSingle) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	defer errs.Wrap(errType, &err, "QueryContext")

	c.printLog(ctx, "QueryContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.conn.QueryContext(ctx, query, args...)
}

func (c *clientSingle) QueryxContext(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	defer errs.Wrap(errType, &err, "QueryxContext")

	c.printLog(ctx, "QueryxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return c.conn.QueryxContext(ctx, query, args...)
}

func (c *clientSingle) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	c.printLog(ctx, "QueryRowxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return c.conn.QueryRowxContext(ctx, query, args...)
}

func (c *clientSingle) PrepareContext(ctx context.Context, query string) (statement *sql.Stmt, err error) {
	defer errs.Wrap(errType, &err, "PrepareContext")

	c.printLog(ctx, "PrepareContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareContext(ctx, query)
	}

	return c.conn.PrepareContext(ctx, query)
}

func (c *clientSingle) NamedExecContext(ctx context.Context, query string, arg interface{}) (result sql.Result, err error) {
	defer errs.Wrap(errType, &err, "NamedExecContext")

	c.printLog(ctx, "NamedExecContext", query, arg)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return c.conn.NamedExecContext(ctx, query, arg)
}

func (c *clientSingle) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	defer errs.Wrap(errType, &err, "SelectContext")

	c.printLog(ctx, "SelectContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return c.conn.SelectContext(ctx, dest, query, args...)
}

func (c *clientSingle) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	defer errs.Wrap(errType, &err, "GetContext")

	c.printLog(ctx, "GetContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return c.conn.GetContext(ctx, dest, query, args...)
}

func (c *clientSingle) PrepareNamedContext(ctx context.Context, query string) (statement *sqlx.NamedStmt, err error) {
	defer errs.Wrap(errType, &err, "PrepareNamedContext")

	c.printLog(ctx, "PrepareNamedContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareNamedContext(ctx, query)
	}

	return c.conn.PrepareNamedContext(ctx, query)
}

func (c *clientSingle) EveryShard(_ func(conn DB) error) error {
	return errors.New("method not supported in single client")
}

func (c *clientSingle) EveryShardAsync(_ func(conn DB) error, _ ...int) error {
	return errors.New("method not supported in single client")
}

func (c *clientSingle) printLog(ctx context.Context, queryType, query string, args ...any) {
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
