package sql

import (
	"fmt"
	"github.com/boostgo/lite/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

// Connect to the database.
// "options" can override default settings
func Connect(connectionString string, options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	connection, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// set default settings
	connection.SetMaxOpenConns(10)
	connection.SetMaxIdleConns(10)
	connection.SetConnMaxLifetime(time.Second * 10)
	connection.SetConnMaxIdleTime(time.Second * 10)

	// apply options
	for _, option := range options {
		option(connection)
	}

	// make ping
	if err = connection.Ping(); err != nil {
		return nil, err
	}

	return connection, nil
}

func MustConnect(connectionString string, options ...func(connection *sqlx.DB)) *sqlx.DB {
	connection, err := Connect(connectionString, options...)
	if err != nil {
		log.Fatal().Err(err).Msg("Connect to Database").Namespace("storage.sql")
	}

	return connection
}

type Connector struct {
	host     string
	port     int
	username string
	password string
	database string

	binaryParameters bool
}

func NewConnector() *Connector {
	return &Connector{}
}

func (connector *Connector) Host(host string) *Connector {
	connector.host = host
	return connector
}

func (connector *Connector) Port(port int) *Connector {
	connector.port = port
	return connector
}

func (connector *Connector) Username(username string) *Connector {
	connector.username = username
	return connector
}

func (connector *Connector) Password(password string) *Connector {
	connector.password = password
	return connector
}

func (connector *Connector) Database(database string) *Connector {
	connector.database = database
	return connector
}

func (connector *Connector) BinaryParameters() *Connector {
	connector.binaryParameters = true
	return connector
}

func (connector *Connector) Build() string {
	var binaryParameters string
	if connector.binaryParameters {
		binaryParameters = " binary_parameters=yes"
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable%s",
		connector.host, connector.port,
		connector.username, connector.password,
		connector.database,
		binaryParameters,
	)
}

func (connector *Connector) String() string {
	return connector.Build()
}

func (connector *Connector) Connect(options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	return Connect(connector.Build(), options...)
}

func (connector *Connector) MustConnect(options ...func(connection *sqlx.DB)) *sqlx.DB {
	return MustConnect(connector.Build(), options...)
}
