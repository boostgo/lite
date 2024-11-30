package postgres

import (
	"context"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/storage/sql"
	"github.com/boostgo/lite/system/health"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

func New(connectionStrings ...string) health.Checker {
	return health.NewChecker("postgres", func(ctx context.Context) (status health.Status, err error) {
		if len(connectionStrings) == 0 {
			return status, errors.New("no postgres connection string provided")
		}

		var wg *errgroup.Group
		wg, ctx = errgroup.WithContext(ctx)
		for _, cs := range connectionStrings {
			wg.Go(func() error {
				return checkConnect(ctx, cs)
			})
		}

		if err = wg.Wait(); err != nil {
			return status, err
		}

		return health.Status{
			Status: health.StatusHealthy,
		}, nil
	})
}

func checkConnect(ctx context.Context, connectionString string) (err error) {
	var conn *sqlx.DB
	conn, err = sql.Connect(connectionString)
	if err != nil {
		return errs.
			New("Health check failed on connecting").
			SetError(err).
			AddContext("connection_string", connectionString)
	}

	defer conn.Close()

	if err = conn.PingContext(ctx); err != nil {
		return errs.
			New("Health check failed on ping").
			SetError(err).
			AddContext("connection_string", connectionString)
	}

	var rows *sqlx.Rows
	rows, err = conn.QueryxContext(ctx, "SELECT VERSION()")
	if err != nil {
		return errs.
			New("Health check failed on SELECT").
			SetError(err).
			AddContext("connection_string", connectionString)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			err = errs.
				New("Health check failed on rows closing").
				SetError(err).
				AddContext("connection_string", connectionString)
			return
		}
	}()

	return nil
}
