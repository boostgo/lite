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

// MustConnect calls Connect and if err catch throws panic
func MustConnect(connectionString string, options ...func(connection *sqlx.DB)) *sqlx.DB {
	connection, err := Connect(connectionString, options...)
	if err != nil {
		log.Fatal().Err(err).Msg("Connect to Database").Namespace("storage.sql")
	}

	return connection
}

// Connector helper for creating connection
type Connector struct {
	host     string
	port     int
	username string
	password string
	database string

	binaryParameters bool
}

// NewConnector creates Connector object
func NewConnector() *Connector {
	return &Connector{}
}

// Host set host of database
func (connector *Connector) Host(host string) *Connector {
	connector.host = host
	return connector
}

// Port set port of database
func (connector *Connector) Port(port int) *Connector {
	connector.port = port
	return connector
}

// Username set username of database user
func (connector *Connector) Username(username string) *Connector {
	connector.username = username
	return connector
}

// Password set password of database user
func (connector *Connector) Password(password string) *Connector {
	connector.password = password
	return connector
}

// Database set database name
func (connector *Connector) Database(database string) *Connector {
	connector.database = database
	return connector
}

// BinaryParameters set binary_parameters=yes param
func (connector *Connector) BinaryParameters() *Connector {
	connector.binaryParameters = true
	return connector
}

// Build connection string
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

// String calls Build method
func (connector *Connector) String() string {
	return connector.Build()
}

// Connect calls Build method and call Connect function
func (connector *Connector) Connect(options ...func(connection *sqlx.DB)) (*sqlx.DB, error) {
	return Connect(connector.Build(), options...)
}

// MustConnect calls MustConnect function
func (connector *Connector) MustConnect(options ...func(connection *sqlx.DB)) *sqlx.DB {
	return MustConnect(connector.Build(), options...)
}
