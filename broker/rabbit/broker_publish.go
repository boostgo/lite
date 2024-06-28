package rabbit

import (
	"context"
	"encoding/json"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/content"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PublishConfig struct {
	Exchange    string
	Mandatory   bool
	Immediate   bool
	ContentType string
}

func defaultPublishConfig() PublishConfig {
	return PublishConfig{
		ContentType: content.JSON,
	}
}

func (broker *Broker) Publish(ctx context.Context, queue string, body any, cfg ...PublishConfig) error {
	bodyInBytes, err := json.Marshal(body)
	if err != nil {
		return errs.New("Parse message body for publish").SetError(err)
	}

	var config PublishConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = defaultPublishConfig()
	}

	headers := amqp.Table{}
	traceIdSet := trace.SetAmqpCtx(ctx, headers)
	if !traceIdSet && broker.isTraceMaster {
		trace.SetAmqp(headers, trace.String())
	}

	if err = broker.channel.PublishWithContext(ctx,
		config.Exchange,
		queue,
		config.Mandatory,
		config.Immediate,
		amqp.Publishing{
			ContentType: config.ContentType,
			Body:        bodyInBytes,
			Headers:     headers,
		}); err != nil {
		return errs.
			New("Publish message error").
			SetError(err).
			AddContext("queue", queue)
	}

	return nil
}
