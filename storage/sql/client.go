package sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

const errType = "SQL"

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
	EveryShard(fn func(conn DB) error) error
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
