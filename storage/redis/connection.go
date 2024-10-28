package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Option func(options *redis.Options)

func Connect(address string, port, db int, password string, opts ...Option) (*redis.Client, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Password: password,
		DB:       db,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func MustConnect(address string, port, db int, password string, opts ...Option) *redis.Client {
	client, err := Connect(address, port, db, password, opts...)
	if err != nil {
		panic(err)
	}

	return client
}
