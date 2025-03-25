package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	"github.com/mailru/go-clickhouse"
	"time"

	"github.com/boostgo/lite/log"
	"github.com/jmoiron/sqlx"
)

func init() {
	// register drivers
	pq.QuoteIdentifier("")    // call package pq to register pq driver
	stdlib.GetDefaultDriver() // call package stdlib to register pgx driver
	clickhouse.Map("")        // call package clickhouse to register ch driver
}

const (
	PqDriver  = "postgres"
	PgxDriver = "pgx"
	ChDriver  = "clickhouse"
)

// RegisterDriver sql package wrap for sql.Register function.
//
// For simple usage of common pkg names (sql)
func RegisterDriver(driverName string, driver driver.Driver) {
	sql.Register(driverName, driver)
}

// Connect to the database.
//
// "options" can override default settings
func Connect(
	driverName, connectionString string,
	timeout time.Duration,
	options ...func(connection *sqlx.DB),
) (*sqlx.DB, error) {
	connection, err := sqlx.Open(driverName, connectionString)
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
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if err = connection.PingContext(ctx); err != nil {
		return nil, err
	}

	return connection, nil
}

// MustConnect calls Connect and if err catch throws panic
func MustConnect(
	driverName, connectionString string,
	timeout time.Duration,
	options ...func(connection *sqlx.DB),
) *sqlx.DB {
	connection, err := Connect(driverName, connectionString, timeout, options...)
	if err != nil {
		log.
			Fatal().
			Namespace("storage.sql").
			Err(err).
			Msg("Connect to Database")
	}

	return connection
}

// MaxConnectionsOption sets max open & idle connections
func MaxConnectionsOption(open, idle int) func(conn *sqlx.DB) {
	return func(conn *sqlx.DB) {
		conn.SetMaxOpenConns(open)
		conn.SetMaxIdleConns(idle)
	}
}

// MaxTimeOption sets connection max lifetime & max idle time settings
func MaxTimeOption(lifetime, idle time.Duration) func(conn *sqlx.DB) {
	return func(conn *sqlx.DB) {
		conn.SetConnMaxLifetime(lifetime)
		conn.SetConnMaxIdleTime(idle)
	}
}

// Connector helper for creating connection
type Connector struct {
	host             string
	port             int
	username         string
	password         string
	database         string
	schema           string
	binaryParameters bool

	timeout time.Duration

	maxOpenConnections int
	maxIdleConnections int
	maxConnLifetime    time.Duration
	maxIdleTime        time.Duration
}

// NewConnector creates Connector object
func NewConnector() *Connector {
	const (
		defaultMaxOpenConnections = 10
		defaultMaxIdleConnections = 10
		defaultMaxConnLifetime    = time.Second * 10
		defaultMaxIdleTime        = time.Second * 10
	)

	return &Connector{
		maxOpenConnections: defaultMaxOpenConnections,
		maxIdleConnections: defaultMaxIdleConnections,
		maxConnLifetime:    defaultMaxConnLifetime,
		maxIdleTime:        defaultMaxIdleTime,
	}
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

// Schema set schema name
func (connector *Connector) Schema(schema string) *Connector {
	connector.schema = schema
	return connector
}

// BinaryParameters set binary_parameters=yes param
func (connector *Connector) BinaryParameters(binaryParameters bool) *Connector {
	connector.binaryParameters = binaryParameters
	return connector
}

// Timeout set timeout for connect & ping (via [context.Context])
func (connector *Connector) Timeout(timeout time.Duration) *Connector {
	connector.timeout = timeout
	return connector
}

// MaxOpenConnections set max open connections option for connection setting
func (connector *Connector) MaxOpenConnections(maxOpenConnections int) *Connector {
	connector.maxOpenConnections = maxOpenConnections
	return connector
}

// MaxIdleConnections set max idle connections option for connection setting
func (connector *Connector) MaxIdleConnections(maxIdleConnections int) *Connector {
	connector.maxIdleConnections = maxIdleConnections
	return connector
}

// MaxIdleTime set max idle time option for connection setting
func (connector *Connector) MaxIdleTime(maxIdleTime time.Duration) *Connector {
	connector.maxIdleTime = maxIdleTime
	return connector
}

// ConnectionMaxLifetime set max connection lifetime option for connection setting
func (connector *Connector) ConnectionMaxLifetime(connectionMaxLifetime time.Duration) *Connector {
	connector.maxConnLifetime = connectionMaxLifetime
	return connector
}

// Build connection string
func (connector *Connector) Build() string {
	var binaryParameters string
	if connector.binaryParameters {
		binaryParameters = " binary_parameters=yes"
	}

	var schema string
	if connector.schema != "" {
		schema = " search_path=" + connector.schema
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable%s%s",
		connector.host, connector.port,
		connector.username, connector.password,
		connector.database,
		binaryParameters,
		schema,
	)
}

// String calls Build method
func (connector *Connector) String() string {
	return connector.Build()
}

// Connect calls Build method and call Connect function
func (connector *Connector) Connect(
	driverName string,
	options ...func(connection *sqlx.DB),
) (*sqlx.DB, error) {
	options = append(
		options,
		MaxConnectionsOption(connector.maxOpenConnections, connector.maxIdleConnections),
		MaxTimeOption(connector.maxConnLifetime, connector.maxIdleTime),
	)

	return Connect(
		driverName,
		connector.Build(),
		connector.timeout,
		options...,
	)
}

// MustConnect calls MustConnect function
func (connector *Connector) MustConnect(
	driverName string,
	options ...func(connection *sqlx.DB),
) *sqlx.DB {
	return MustConnect(
		driverName,
		connector.Build(),
		connector.timeout,
		options...,
	)
}

// ConnectionString returns built connection string by provided params for sqlx lib
func ConnectionString(
	host string,
	port int,
	username, password string,
	database string,
	binaryParameters bool,
) string {
	return NewConnector().
		Host(host).
		Port(port).
		Username(username).
		Password(password).
		Database(database).
		BinaryParameters(binaryParameters).
		String()
}
