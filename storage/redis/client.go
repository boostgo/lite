package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/list"
	"github.com/redis/go-redis/v9"
	"time"
)

const errType = "Redis"

var (
	ErrKeyEmpty = errors.New("key is empty")
)

type Client struct {
	client *redis.Client
}

func New(address string, port, db int, password string, opts ...Option) (*Client, error) {
	conn, err := Connect(address, port, db, password, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: conn,
	}, nil
}

func Must(address string, port, db int, password string, opts ...Option) *Client {
	client, err := New(address, port, db, password, opts...)
	if err != nil {
		panic(err)
	}

	return client
}

func NewFromClient(conn *redis.Client) *Client {
	return &Client{
		client: conn,
	}
}

func (client *Client) Close() error {
	return client.client.Close()
}

func (client *Client) Client() *redis.Client {
	return client.client
}

func (client *Client) Pipeline() redis.Pipeliner {
	return client.client.Pipeline()
}

func (client *Client) TxPipeline() redis.Pipeliner {
	return client.client.TxPipeline()
}

func (client *Client) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "Keys")
	return client.client.Keys(ctx, pattern).Result()
}

func (client *Client) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return nil
	}

	defer errs.Wrap(errType, &err, "Delete")

	// clean up keys from empty
	keys = list.Filter(keys, func(key string) bool {
		return key != ""
	})

	if len(keys) == 0 {
		return nil
	}

	return client.client.Del(ctx, keys...).Err()
}

func (client *Client) Dump(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Dump")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.Dump(ctx, key).Result()
}

func (client *Client) Rename(ctx context.Context, oldKey, newKey string) (err error) {
	defer errs.Wrap(errType, &err, "Rename")

	if err = validateKey(oldKey); err != nil {
		return errs.
			New("Old key is invalid").
			SetError(err)
	}

	if err = validateKey(newKey); err != nil {
		return errs.
			New("New key is invalid").
			SetError(err)
	}

	return client.client.Rename(ctx, oldKey, newKey).Err()
}

func (client *Client) Refresh(ctx context.Context, key string, ttl time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Refresh")

	if err = validateKey(key); err != nil {
		return err
	}

	return client.client.Expire(ctx, key, ttl).Err()
}

func (client *Client) RefreshAt(ctx context.Context, key string, at time.Time) (err error) {
	defer errs.Wrap(errType, &err, "RefreshAt")

	if err = validateKey(key); err != nil {
		return err
	}

	return client.client.ExpireAt(ctx, key, at).Err()
}

func (client *Client) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	defer errs.Wrap(errType, &err, "TTL")

	if err = validateKey(key); err != nil {
		return ttl, err
	}

	ttl, err = client.client.TTL(ctx, key).Result()
	if err != nil {
		return ttl, err
	}

	const notExistKey = -2
	if ttl == notExistKey {
		return ttl, errs.ErrNotFound
	}

	return ttl, nil
}

func (client *Client) Set(ctx context.Context, key string, value any, ttl ...time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Set")

	if err = validateKey(key); err != nil {
		return err
	}

	var expireAt time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = ttl[0]
	}

	return client.client.Set(ctx, key, value, expireAt).Err()
}

func (client *Client) Get(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	result, err = client.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, errs.
				New("Redis key not found").
				SetError(errs.ErrNotFound).
				AddContext("key", key)
		}

		return result, err
	}

	return result, nil
}

func (client *Client) Exist(ctx context.Context, key string) (result int64, err error) {
	defer errs.Wrap(errType, &err, "Exist")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.Exists(ctx, key).Result()
}

func (client *Client) GetBytes(ctx context.Context, key string) (result []byte, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	result, err = client.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, errs.
				New("Redis key not found").
				SetError(errs.ErrNotFound).
				AddContext("key", key)
		}

		return result, err
	}

	return result, nil
}

func (client *Client) GetInt(ctx context.Context, key string) (result int, err error) {
	defer errs.Wrap(errType, &err, "GetInt")

	if err = validateKey(key); err != nil {
		return result, err
	}

	result, err = client.client.Get(ctx, key).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, errs.
				New("Redis key not found").
				SetError(errs.ErrNotFound).
				AddContext("key", key)
		}

		return result, err
	}

	return result, nil
}

func (client *Client) Parse(ctx context.Context, key string, export any) (err error) {
	defer errs.Wrap(errType, &err, "Parse")

	if err = validateKey(key); err != nil {
		return err
	}

	var result []byte
	result, err = client.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errs.
				New("Redis key not found").
				SetError(errs.ErrNotFound).
				AddContext("key", key)
		}

		return err
	}

	return json.Unmarshal(result, &export)
}

func (client *Client) HSet(ctx context.Context, key string, value map[string]any) (err error) {
	defer errs.Wrap(errType, &err, "HSet")

	if err = validateKey(key); err != nil {
		return err
	}

	return client.client.HSet(ctx, key, value).Err()
}

func (client *Client) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	defer errs.Wrap(errType, &err, "HGetAll")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.HGetAll(ctx, key).Result()
}

func (client *Client) HGet(ctx context.Context, key, field string) (result string, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.HGet(ctx, key, field).Result()
}

func (client *Client) HGetInt(ctx context.Context, key, field string) (result int, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.HGet(ctx, key, field).Int()
}

func (client *Client) HGetBool(ctx context.Context, key, field string) (result bool, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	return client.client.HGet(ctx, key, field).Bool()
}

func (client *Client) HExist(ctx context.Context, key, field string) (exist bool, err error) {
	defer errs.Wrap(errType, &err, "HExist")

	if err = validateKey(key); err != nil {
		return exist, err
	}

	return client.client.HExists(ctx, key, field).Result()
}

func (client *Client) Scan(ctx context.Context, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "Scan")
	return client.client.Scan(ctx, cursor, pattern, count).Result()
}

func validateKey(key string) error {
	if key == "" {
		return ErrKeyEmpty
	}

	return nil
}
