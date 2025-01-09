package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

const errType = "Redis"

var (
	ErrKeyEmpty = errors.New("key is empty")
)

type Client interface {
	io.Closer

	Client(ctx context.Context) (redis.UniversalClient, error)
	Pipeline(ctx context.Context) (redis.Pipeliner, error)
	TxPipeline(ctx context.Context) (redis.Pipeliner, error)

	Keys(ctx context.Context, pattern string) ([]string, error)
	Delete(ctx context.Context, keys ...string) error
	Dump(ctx context.Context, key string) (string, error)
	Rename(ctx context.Context, oldKey, newKey string) error
	Refresh(ctx context.Context, key string, ttl time.Duration) error
	RefreshAt(ctx context.Context, key string, at time.Time) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Set(ctx context.Context, key string, value any, ttl ...time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Exist(ctx context.Context, key string) (int64, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	GetInt(ctx context.Context, key string) (int, error)
	Parse(ctx context.Context, key string, export any) error
	Scan(ctx context.Context, cursor uint64, pattern string, count int64) ([]string, uint64, error)

	HSet(ctx context.Context, key string, value map[string]any) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetInt(ctx context.Context, key, field string) (int, error)
	HGetBool(ctx context.Context, key, field string) (bool, error)
	HExist(ctx context.Context, key, field string) (bool, error)
}

func validateKey(key string) error {
	if key == "" {
		return ErrKeyEmpty
	}

	return nil
}
