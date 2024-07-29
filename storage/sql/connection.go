package sql

import (
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

func Connect(connectionString string, options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	connection, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	connection.SetMaxOpenConns(10)
	connection.SetMaxIdleConns(10)
	connection.SetConnMaxLifetime(time.Second * 10)
	connection.SetConnMaxIdleTime(time.Second * 10)

	for _, option := range options {
		option(connection)
	}

	if err = connection.Ping(); err != nil {
		return nil, err
	}

	life.Tear(connection.Close)

	return connection, nil
}

func MustConnect(connectionString string, options ...func(connection *sqlx.DB)) *sqlx.DB {
	connection, err := Connect(connectionString, options...)
	if err != nil {
		log.Fatal().Err(err).Msg("Connect to Database").Namespace("storage.sql")
	}

	return connection
}
