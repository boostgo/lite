package bus

import (
	"context"
	"encoding/json"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/trace"
	amqp "github.com/rabbitmq/amqp091-go"
)

type eventContext struct {
	body   []byte
	ctx    context.Context
	cancel func()
	logger log.Logger
}

func newContextRMQ(ctx context.Context, message *amqp.Delivery) EventContext {
	return &eventContext{
		ctx:    ctx,
		body:   message.Body,
		logger: log.Namespace("event-context"),
	}
}

func (ctx *eventContext) Body() []byte {
	return ctx.body
}

func (ctx *eventContext) Parse(object any) error {
	return json.Unmarshal(ctx.body, &object)
}

func (ctx *eventContext) Context() context.Context {
	if ctx.ctx == nil {
		return context.Background()
	}

	return ctx.ctx
}

func (ctx *eventContext) TraceID() string {
	return trace.Get(ctx.ctx)
}

func (ctx *eventContext) Cancel() {
	if ctx.cancel == nil {
		return
	}

	ctx.cancel()
}

func (ctx *eventContext) Logger() log.Logger {
	return ctx.logger
}
