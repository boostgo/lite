package connector

import (
	"github.com/boostgo/lite/collections/concurrent"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type Connector struct {
	connections *concurrent.Map[string, *amqp.Connection]
}

var (
	_connector *Connector
	_once      sync.Once
)

func Get() *Connector {
	_once.Do(func() {
		_connector = &Connector{
			connections: concurrent.NewMap[string, *amqp.Connection](),
		}
	})

	return _connector
}

func (connector *Connector) Register(connectionString string) error {
	_, ok := connector.connections.Load(connectionString)
	if ok {
		return nil
	}

	conn, err := connector.connect(connectionString)
	if err != nil {
		return err
	}

	connector.connections.Store(connectionString, conn)
	return nil
}

func (connector *Connector) Get(connectionString string) (*amqp.Connection, error) {
	conn, ok := connector.connections.Load(connectionString)
	if !ok {
		if err := connector.Register(connectionString); err != nil {
			return nil, err
		}

		conn, _ = connector.connections.Load(connectionString)
		return conn, nil
	}

	return conn, nil
}

func (connector *Connector) Close() error {
	connector.connections.Each(func(_ string, value *amqp.Connection) bool {
		_ = value.Close()
		return true
	})

	for _, key := range connector.connections.Keys() {
		connector.connections.Delete(key)
	}

	return nil
}

func (connector *Connector) connect(connectionString string) (*amqp.Connection, error) {
	return amqp.Dial(connectionString)
}
