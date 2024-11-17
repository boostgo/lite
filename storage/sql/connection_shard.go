package sql

import (
	"context"
	"database/sql"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/storage"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type ShardConnectString struct {
	Key              string
	ConnectionString string
}

type ShardConnect interface {
	Key() string
	Conn() *sqlx.DB
	Close() error
}

type shardConnect struct {
	key  string
	conn *sqlx.DB
}

func newShardConnect(key string, conn *sqlx.DB) ShardConnect {
	return &shardConnect{
		key:  key,
		conn: conn,
	}
}

func (conn *shardConnect) Key() string {
	return conn.key
}

func (conn *shardConnect) Conn() *sqlx.DB {
	return conn.conn
}

func (conn *shardConnect) Close() error {
	return conn.conn.Close()
}

func ConnectShards(connectionStrings []ShardConnectString, selector ConnectionSelector, options ...func(connection *sqlx.DB)) (*Connections, error) {
	// validate for connection key unique and for empty
	// also, validate for empty connection string
	keys := make(map[string]struct{}, len(connectionStrings))
	for _, cs := range connectionStrings {
		if cs.Key == "" {
			return nil, errs.New("Connection key is empty")
		}

		if cs.ConnectionString == "" {
			return nil, errs.
				New("Connection string is empty").
				AddContext("key", cs.Key)
		}

		if _, ok := keys[cs.Key]; ok {
			return nil, errs.
				New("Connection keys cannot duplicate").
				AddContext("key", cs.Key)
		}

		keys[cs.Key] = struct{}{}
	}

	// connect every shard
	connections := make([]ShardConnect, len(connectionStrings))
	for idx, cs := range connectionStrings {
		connection, err := Connect(cs.ConnectionString, options...)
		if err != nil {
			return nil, err
		}

		connections[idx] = newShardConnect(cs.Key, connection)
	}

	return newConnections(connections, selector), nil
}

func MustConnectShards(connectionStrings []ShardConnectString, selector ConnectionSelector, options ...func(connection *sqlx.DB)) *Connections {
	connections, err := ConnectShards(connectionStrings, selector, options...)
	if err != nil {
		panic(err)
	}

	return connections
}

type Connections struct {
	connections []ShardConnect
	selector    ConnectionSelector
}

func newConnections(connections []ShardConnect, selector ConnectionSelector) *Connections {
	return &Connections{
		connections: connections,
		selector:    selector,
	}
}

func (c *Connections) Get(ctx context.Context) (ShardConnect, error) {
	// get shard by provided selector
	conn := c.selector(ctx, c.connections)
	if conn == nil {
		return nil, storage.ErrConnNotSelected
	}

	return conn, nil
}

func (c *Connections) Connections() []ShardConnect {
	return c.connections
}

func (c *Connections) RawConnections() []*sqlx.DB {
	connections := make([]*sqlx.DB, len(c.connections))
	for idx, conn := range c.connections {
		connections[idx] = conn.Conn()
	}
	return connections
}

func (c *Connections) Close() error {
	wg := errgroup.Group{}

	for _, conn := range c.connections {
		wg.Go(conn.Close)
	}

	return wg.Wait()
}

func (c *Connections) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	// begin transaction at selected shard
	conn, err := c.Get(ctx)
	if err != nil {
		return nil, err
	}

	return conn.Conn().BeginTxx(ctx, opts)
}
