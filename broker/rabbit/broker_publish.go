package rabbit

import (
	"context"
	"encoding/json"

	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/system/trace"
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
		ContentType: "application/json",
	}
}

// Publish send provided body to queue and optional configurations.
//
// If context has trace id it provides to header of sending message
func (broker *Broker) Publish(ctx context.Context, queue string, body any, cfg ...PublishConfig) error {
	bodyInBytes, err := json.Marshal(body)
	if err != nil {
		return errorx.
			New("Parse message body for publish").
			SetError(err)
	}

	var config PublishConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = defaultPublishConfig()
	}

	headers := amqp.Table{}
	if broker.isTraceMaster {
		if trace.Get(ctx) == "" {
			ctx = trace.Set(ctx, trace.String())
		}

		trace.SetAmqp(headers, ctx)
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
		return errorx.
			New("Publish message error").
			SetError(err).
			AddContext("queue", queue)
	}

	return nil
}
