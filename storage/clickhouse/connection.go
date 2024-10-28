package clickhouse

import (
	"context"
	"github.com/boostgo/lite/log"
	"github.com/jmoiron/sqlx"
	"time"

	_ "github.com/mailru/go-clickhouse"
)

func Connect(connectionString string, options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	connection, err := sqlx.Connect("clickhouse", connectionString)
	if err != nil {
		return nil, err
	}

	connection.SetMaxOpenConns(20)
	connection.SetMaxIdleConns(5)
	connection.SetConnMaxIdleTime(time.Second * 30)
	connection.SetConnMaxLifetime(time.Second * 60)

	for _, option := range options {
		option(connection)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err = connection.PingContext(ctx); err != nil {
		return nil, err
	}

	return connection, nil
}

func MustConnect(connectionString string, options ...func(connection *sqlx.DB)) *sqlx.DB {
	connection, err := Connect(connectionString, options...)
	if err != nil {
		log.Fatal().Err(err).Msg("Connect to Clickhouse")
	}

	return connection
}
