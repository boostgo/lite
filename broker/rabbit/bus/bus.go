package bus

import (
	"context"
	"github.com/boostgo/lite/broker/rabbit/connector"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/validator"
	"reflect"
	"sync"
	"time"
)

// Bus is manager of event bus (or message bus) pattern.
//
// It creates Dispatcher and Listener objects to work with this pattern.
//
// How it works:
//  1. Creates exchanges: "message.bus" and "message.bus.errors"
//  2. Bindings in Listener to queues their consuming actions, it creates need queues automatically
//  3. When dispatcher sends event/message with some structure, it takes structure name and send event/message to queue named as structure name
type Bus struct {
	connector *connector.Connector
	validator *validator.Validator
}

const (
	namespace = "bus"

	defaultMessageBusExchangeName       = "message.bus"
	defaultMessageBusErrorsExchangeName = "message.bus.errors"

	defaultTimeout = time.Minute
)

var (
	_once sync.Once
	_bus  *Bus
)

// Get creates Bus object if it does not initialize.
//
// If instance initialized - returns singleton object
func Get() *Bus {
	_once.Do(func() {
		v, err := validator.New()
		if err != nil {
			log.Namespace(namespace).Fatal().Err(err).Msg("Create validator")
		}

		_bus = &Bus{
			connector: connector.Get(),
			validator: v,
		}
	})

	return _bus
}

// Dispatcher creates Dispatcher with provided connection string.
//
// It creates new connection & channel for its own
func (bus *Bus) Dispatcher(ctx context.Context, connectionString string) (Dispatcher, error) {
	conn, err := bus.connector.Get(connectionString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return newDispatcher(ctx, channel, bus.validator), nil
}

// Listener creates Listener with provided connection string.
//
// It creates new connection & channel for its own
func (bus *Bus) Listener(ctx context.Context, connectionString string) (Listener, error) {
	conn, err := bus.connector.Get(connectionString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return newListener(ctx, channel, bus.validator), nil
}

func nameOfEvent(event any) (string, error) {
	reflectValue := reflect.TypeOf(event)
	if reflectValue.Kind() != reflect.Struct {
		return "", errs.
			New("Provided argument is not struct").
			AddContext("actual-type", reflectValue.Kind())
	}

	return reflectValue.Name(), nil
}
