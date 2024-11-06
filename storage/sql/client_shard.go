package sql

import (
	"context"
	"database/sql"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/storage"
	"github.com/jmoiron/sqlx"
)

type ConnectionSelector func(ctx context.Context, connections []ShardConnect) ShardConnect

type clientShard struct {
	connections *Connections
	enableLog   bool
}

func ClientShard(connections *Connections, enableLog ...bool) DB {
	var enable bool
	if len(enableLog) > 0 {
		enable = enableLog[0]
	}

	return &clientShard{
		connections: connections,
		enableLog:   enable,
	}
}

func (c *clientShard) Connection() *sqlx.DB {
	return nil
}

func (c *clientShard) ExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	defer errs.Wrap(errType, &err, "ExecContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}
	c.printLog(ctx, raw.Name(), "ExecContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return raw.Conn().ExecContext(ctx, query, args...)
}

func (c *clientShard) QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	defer errs.Wrap(errType, &err, "QueryContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}
	c.printLog(ctx, raw.Name(), "QueryContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return raw.Conn().QueryContext(ctx, query, args...)
}

func (c *clientShard) QueryxContext(ctx context.Context, query string, args ...interface{}) (rows *sqlx.Rows, err error) {
	defer errs.Wrap(errType, &err, "QueryxContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}
	c.printLog(ctx, raw.Name(), "QueryxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return raw.Conn().QueryxContext(ctx, query, args...)
}

func (c *clientShard) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil
	}

	c.printLog(ctx, raw.Name(), "QueryRowxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return raw.Conn().QueryRowxContext(ctx, query, args...)
}

func (c *clientShard) PrepareContext(ctx context.Context, query string) (statement *sql.Stmt, err error) {
	defer errs.Wrap(errType, &err, "PrepareContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}
	c.printLog(ctx, raw.Name(), "PrepareContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareContext(ctx, query)
	}

	return raw.Conn().PrepareContext(ctx, query)
}

func (c *clientShard) NamedExecContext(ctx context.Context, query string, arg interface{}) (result sql.Result, err error) {
	defer errs.Wrap(errType, &err, "NamedExecContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}
	c.printLog(ctx, raw.Name(), "NamedExecContext", query, arg)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return raw.Conn().NamedExecContext(ctx, query, arg)
}

func (c *clientShard) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	defer errs.Wrap(errType, &err, "SelectContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return err
	}
	c.printLog(ctx, raw.Name(), "SelectContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return raw.Conn().SelectContext(ctx, dest, query, args...)
}

func (c *clientShard) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	defer errs.Wrap(errType, &err, "GetContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return err
	}
	c.printLog(ctx, raw.Name(), "GetContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return raw.Conn().GetContext(ctx, dest, query, args...)
}

func (c *clientShard) PrepareNamedContext(ctx context.Context, query string) (statement *sqlx.NamedStmt, err error) {
	defer errs.Wrap(errType, &err, "PrepareNamedContext")

	raw, err := c.selectConnect(ctx)
	if err != nil {
		return nil, err
	}

	c.printLog(ctx, raw.Name(), "PrepareNamedContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareNamedContext(ctx, query)
	}

	return raw.Conn().PrepareNamedContext(ctx, query)
}

func (c *clientShard) printLog(ctx context.Context, connectionName, queryType, query string, args ...any) {
	if !c.enableLog || storage.IsNoLog(ctx) {
		return
	}

	log.
		Context(ctx, "storage.sql."+queryType).
		Info().
		Str("connection_name", connectionName).
		Str("query", query).
		Any("args", args).
		Send()
}

func (c *clientShard) selectConnect(ctx context.Context) (ShardConnect, error) {
	return c.connections.Get(ctx)
}