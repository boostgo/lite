package sql

import (
	"context"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func Migrate(ctx context.Context, conn *sqlx.DB, databaseName string) (err error) {
	const errType = "Storage Migrate"
	defer errs.Wrap(errType, &err, "Migrate")

	nativeConn, err := conn.Conn(ctx)
	if err != nil {
		return err
	}

	driver, err := postgres.WithConnection(ctx, nativeConn, &postgres.Config{})
	if err != nil {
		return err
	}

	_, err = nativeConn.ExecContext(ctx, "SET lock_timeout = '60s';")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://./migrations", databaseName, driver)
	if err != nil {
		return err
	}
	defer migrator.Close()

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.
				Info().
				Str("database_name", databaseName).
				Msg("Migrate no changes")
			return nil
		}

		return err
	}

	return nil
}

func MustMigrate(ctx context.Context, conn *sqlx.DB, databaseName string) {
	if err := Migrate(ctx, conn, databaseName); err != nil {
		panic(err)
	}
}

func BackgroundMigrate(ctx context.Context, conn *sqlx.DB, databaseName string) {
	_ = Migrate(ctx, conn, databaseName)
}
