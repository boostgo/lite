package sql

import (
	"github.com/boostgo/lite/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(connectionString string, options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	connection, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		option(connection)
	}

	if err = connection.Ping(); err != nil {
		return nil, err
	}

	return connection, nil
}

func MustConnect(connectionString string, options ...func(connection *sqlx.DB)) *sqlx.DB {
	connection, err := Connect(connectionString, options...)
	if err != nil {
		log.Fatal("storage.sql").Err(err).Msg("Connect to Database")
	}

	return connection
}
