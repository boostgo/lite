package mongo

import (
	"context"
	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/storage/mongo"
	"github.com/boostgo/lite/system/health"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

func New(username, password, host string, port int) health.Checker {
	return health.NewChecker("mongo", func(ctx context.Context) (status health.Status, err error) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, username, password, host, port)
		if err != nil {
			return status, errorx.
				New("Health check failed on connect").
				SetError(err).
				AddContext("host", host).
				AddContext("port", port).
				AddContext("username", username)
		}

		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			return status, errorx.
				New("Health check failed on ping").
				SetError(err).
				AddContext("host", host).
				AddContext("port", port).
				AddContext("username", username)
		}

		return health.Status{
			Status: health.StatusHealthy,
		}, nil
	})
}
