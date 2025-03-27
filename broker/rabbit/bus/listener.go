package bus

import (
	"context"
	"github.com/boostgo/appx"
	"time"

	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/broker/rabbit"
	"github.com/boostgo/lite/broker/rabbit/exchanges"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/validator"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ListenerConfig struct {
	MessageBusExchangeName       string
	MessageBusErrorsExchangeName string
	Timeout                      time.Duration
}

func defaultListenerConfig() ListenerConfig {
	return ListenerConfig{
		MessageBusExchangeName:       defaultMessageBusExchangeName,
		MessageBusErrorsExchangeName: defaultMessageBusErrorsExchangeName,
		Timeout:                      defaultTimeout,
	}
}

type listener struct {
	logger log.Logger

	ctx       context.Context
	broker    *rabbit.Broker
	validator *validator.Validator
	config    ListenerConfig

	events []Event
}

func newListener(ctx context.Context, channel *amqp.Channel, validator *validator.Validator, cfg ...ListenerConfig) *listener {
	var listenerConfig ListenerConfig
	if len(cfg) > 0 {
		listenerConfig = cfg[0]
	} else {
		listenerConfig = defaultListenerConfig()
	}

	return &listener{
		logger: log.Namespace("bus.listener"),

		ctx:       ctx,
		broker:    rabbit.NewBroker(channel),
		validator: validator,
		config:    listenerConfig,

		events: make([]Event, 0, 10),
	}
}

func (l *listener) Run() error {
	if err := l.declareExchanges(); err != nil {
		return err
	}

	if err := l.declareEvents(); err != nil {
		return err
	}

	for _, event := range l.events {
		messages, err := l.broker.Consume(event.Name)
		if err != nil {
			return err
		}

		l.listenQueue(messages, event)
	}

	appx.Tear(l.Close)
	return nil
}

func (l *listener) Bind(event any, action func(ctx EventContext) error) {
	eventName, err := nameOfEvent(event)
	if err != nil {
		return
	}

	l.events = append(l.events, Event{
		Name:   eventName,
		Action: action,
		Object: event,
	})
}

func (l *listener) EventsCount() int {
	return len(l.events)
}

func (l *listener) Close() error {
	return l.broker.Close()
}

func (l *listener) declareExchanges() error {
	// message.bus
	if err := l.broker.NewExchange(l.config.MessageBusExchangeName, exchanges.Direct); err != nil {
		return err
	}

	// message.bus.errors
	if err := l.broker.NewExchange(l.config.MessageBusErrorsExchangeName, exchanges.Direct); err != nil {
		return err
	}

	return nil
}

func (l *listener) declareEvents() error {
	for _, event := range l.events {
		// first, create queue
		if _, err := l.broker.NewQueue(event.Name); err != nil {
			return err
		}

		// after this, bind exchange & queue
		if err := l.broker.Bind(l.config.MessageBusExchangeName, event.Name); err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) listenQueue(queue rabbit.MessagesQueue, event Event) {
	go func() {
		for {
			message, ok := <-queue
			if !ok {
				return
			}

			if err := errorx.Try(func() error {
				return l.callEventAction(&message, &event)
			}); err != nil {
				l.logger.Error().Err(err).Msg("Event action")
			}
		}
	}()
}

func (l *listener) callEventAction(message *amqp.Delivery, event *Event) error {
	var ctx context.Context
	var ctxCancel func()
	if l.config.Timeout != 0 {
		ctx, ctxCancel = context.WithTimeout(context.Background(), l.config.Timeout)
	} else {
		ctx = context.Background()
		ctxCancel = func() {}
	}
	defer ctxCancel()

	traceID := trace.GetAMQP(message)
	if traceID != "" {
		ctx = trace.Set(ctx, traceID)
	}

	eventCtx := newContextRMQ(ctx, message)
	defer eventCtx.Cancel()

	if err := event.Action(eventCtx); err != nil {
		return err
	}

	return l.broker.Ack(message.DeliveryTag, false)
}
