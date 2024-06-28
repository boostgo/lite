package bus

import (
	"context"
	"github.com/boostgo/lite/broker/rabbit"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/validator"
	amqp "github.com/rabbitmq/amqp091-go"
)

type dispatcher struct {
	logger log.Logger

	ctx       context.Context
	broker    *rabbit.Broker
	validator *validator.Validator
}

func newDispatcher(ctx context.Context, channel *amqp.Channel, validator *validator.Validator) *dispatcher {
	d := &dispatcher{
		logger: log.Namespace("bus.dispatcher"),

		ctx:       ctx,
		broker:    rabbit.NewBroker(channel),
		validator: validator,
	}

	life.Tear(d.Close)
	return d
}

func (d *dispatcher) Dispatch(ctx context.Context, event any) (err error) {
	if d.ctx.Err() != nil {
		return nil
	}

	defer func() {
		if err == nil {
			return
		}

		d.logger.Error().Err(err).Msg("Dispatch")
	}()

	eventName, err := nameOfEvent(event)
	if err != nil {
		return err
	}

	if err = d.validator.Struct(event); err != nil {
		return err
	}

	if err = d.broker.Publish(ctx, eventName, event); err != nil {
		return errs.New("Dispatch message").SetError(err)
	}

	return nil
}

func (d *dispatcher) Close() error {
	return d.broker.Close()
}
