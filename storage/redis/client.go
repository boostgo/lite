package redis

import (
	"context"
	"errors"
	"github.com/boostgo/lite/list"
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
	MGet(ctx context.Context, keys []string) ([]any, error)
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
	HDelete(ctx context.Context, key string, fields ...string) error
	HScan(ctx context.Context, key string, cursor uint64, pattern string, count int64) ([]string, uint64, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HLen(ctx context.Context, key string) (int64, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]any, error)
	HMSet(ctx context.Context, key string, values ...any) error
	HSetNX(ctx context.Context, key, field string, value any) error
	HScanNoValues(ctx context.Context, key string, cursor uint64, pattern string, count int64) ([]string, uint64, error)
	HVals(ctx context.Context, key string) ([]string, error)
	HRandField(ctx context.Context, key string, count int) ([]string, error)
	HRandFieldWithValues(ctx context.Context, key string, count int) ([]redis.KeyValue, error)
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) ([]int64, error)
	HTTL(ctx context.Context, key string, fields ...string) ([]int64, error)
}

func validateKey(key string) error {
	if key == "" {
		return ErrKeyEmpty
	}

	return nil
}

func validateKeys(keys []string) {
	keys = list.Filter(keys, func(s string) bool {
		return validateKey(s) != nil
	})
}
