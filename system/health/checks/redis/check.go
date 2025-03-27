package redis

import (
	"context"
	"errors"
	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/storage/redis"
	"github.com/boostgo/lite/system/health"
	"golang.org/x/sync/errgroup"
)

func New(cfg ...redis.ConnectionConfig) health.Checker {
	return health.NewChecker("redis", func(ctx context.Context) (status health.Status, err error) {
		if len(cfg) == 0 {
			return status, errors.New("no redis connection string provided")
		}

		var wg *errgroup.Group
		wg, ctx = errgroup.WithContext(ctx)
		for _, c := range cfg {
			wg.Go(func() error {
				return checkConnect(ctx, c)
			})
		}

		if err = wg.Wait(); err != nil {
			return status, err
		}

		return health.Status{
			Status: health.StatusHealthy,
		}, nil
	})
}

func checkConnect(ctx context.Context, cfg redis.ConnectionConfig) (err error) {
	client, err := redis.Connect(cfg.Address, cfg.Port, cfg.DB, cfg.Password)
	if err != nil {
		return errorx.
			New("Health check failed on connect client").
			SetError(err)
	}
	defer client.Close()

	result, err := client.Ping(ctx).Result()
	if err != nil {
		return errorx.
			New("Health check failed on ping/pong").
			SetError(err)
	}

	if result != "PONG" {
		return errorx.
			New("Health check failed on ping/pong result compare").
			AddContext("result", result)
	}

	return nil
}
