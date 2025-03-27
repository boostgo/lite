package mongo

import (
	"context"
	"fmt"
	"github.com/boostgo/lite/system/life"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

const (
	BuildTest = "test"
)

func Connect(ctx context.Context, username, password, host string, port int, opts ...options.Lister[options.ClientOptions]) (*mongo.Client, error) {
	connectionOpt := options.
		Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port))

	client, err := mongo.Connect(
		append(
			[]options.Lister[options.ClientOptions]{
				connectionOpt,
			},
			opts...,
		)...,
	)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func Must(ctx context.Context, username, password, host string, port int, opts ...options.Lister[options.ClientOptions]) *mongo.Client {
	client, err := Connect(ctx, username, password, host, port, opts...)
	if err != nil {
		panic(err)
	}

	life.Tear(func() error {
		tearCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return client.Disconnect(tearCtx)
	})
	return client
}
