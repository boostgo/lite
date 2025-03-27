package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/boostgo/errorx"
	"github.com/jmoiron/sqlx"
)

const errType = "SQL"

// NotFound check if provided error is not found error
func NotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, errorx.ErrNotFound)
}

// DB description of all methods of sqlx package.
//
// Can be used as single client & shard client
type DB interface {
	Connection() *sqlx.DB
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
	GetContext
	NamedExecContext
	SelectContext
	PrepareContext
	EachShard(fn func(conn DB) error) error
	EachShardAsync(fn func(conn DB) error, limit ...int) error
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
