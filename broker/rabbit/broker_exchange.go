package rabbit

import (
	"github.com/boostgo/lite/errs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeConfig struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

func defaultExchangeConfig() ExchangeConfig {
	return ExchangeConfig{}
}

// NewExchange create new exchange with provided type and optional configurations
func (broker *Broker) NewExchange(name, exchangeType string, cfg ...ExchangeConfig) error {
	var config ExchangeConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = defaultExchangeConfig()
	}

	if err := broker.channel.ExchangeDeclare(
		name,
		exchangeType,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	); err != nil {
		return errs.New("Create new exchange error").SetError(err).SetContext(map[string]any{
			"exchange-name": name,
			"exchange-type": exchangeType,
		})
	}

	return nil
}
