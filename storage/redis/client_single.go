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

type singleClient struct {
	client redis.UniversalClient
}

func New(address string, port, db int, password string, opts ...Option) (Client, error) {
	conn, err := Connect(address, port, db, password, opts...)
	if err != nil {
		return nil, err
	}

	return &singleClient{
		client: conn,
	}, nil
}

func Must(address string, port, db int, password string, opts ...Option) Client {
	client, err := New(address, port, db, password, opts...)
	if err != nil {
		panic(err)
	}

	return client
}

func NewFromClient(conn *redis.Client) Client {
	return &singleClient{
		client: conn,
	}
}

func (client *singleClient) Close() error {
	return client.client.Close()
}

func (client *singleClient) Client(_ context.Context) (redis.UniversalClient, error) {
	return client.client, nil
}

func (client *singleClient) Pipeline(_ context.Context) (redis.Pipeliner, error) {
	return client.client.Pipeline(), nil
}

func (client *singleClient) TxPipeline(_ context.Context) (redis.Pipeliner, error) {
	return client.client.TxPipeline(), nil
}

func (client *singleClient) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "Keys")

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Keys(ctx, pattern).Result()
	}

	return client.client.Keys(ctx, pattern).Result()
}

func (client *singleClient) Delete(ctx context.Context, keys ...string) (err error) {
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

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Del(ctx, keys...).Err()
	}

	return client.client.Del(ctx, keys...).Err()
}

func (client *singleClient) Dump(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Dump")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Dump(ctx, key).Result()
	}

	return client.client.Dump(ctx, key).Result()
}

func (client *singleClient) Rename(ctx context.Context, oldKey, newKey string) (err error) {
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

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Rename(ctx, oldKey, newKey).Err()
	}

	return client.client.Rename(ctx, oldKey, newKey).Err()
}

func (client *singleClient) Refresh(ctx context.Context, key string, ttl time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Refresh")

	if err = validateKey(key); err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Expire(ctx, key, ttl).Err()
	}

	return client.client.Expire(ctx, key, ttl).Err()
}

func (client *singleClient) RefreshAt(ctx context.Context, key string, at time.Time) (err error) {
	defer errs.Wrap(errType, &err, "RefreshAt")

	if err = validateKey(key); err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExpireAt(ctx, key, at).Err()
	}

	return client.client.ExpireAt(ctx, key, at).Err()
}

func (client *singleClient) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	defer errs.Wrap(errType, &err, "TTL")

	if err = validateKey(key); err != nil {
		return ttl, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		ttl, err = tx.TTL(ctx, key).Result()
	} else {
		ttl, err = client.client.TTL(ctx, key).Result()
	}
	if err != nil {
		return ttl, err
	}

	const notExistKey = -2
	if ttl == notExistKey {
		return ttl, errs.ErrNotFound
	}

	return ttl, nil
}

func (client *singleClient) Set(ctx context.Context, key string, value any, ttl ...time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Set")

	if err = validateKey(key); err != nil {
		return err
	}

	var expireAt time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = ttl[0]
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Set(ctx, key, value, expireAt).Err()
	}

	return client.client.Set(ctx, key, value, expireAt).Err()
}

func (client *singleClient) Get(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Result()
	} else {
		result, err = client.client.Get(ctx, key).Result()
	}
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

func (client *singleClient) MGet(ctx context.Context, keys []string) (result []any, err error) {
	defer errs.Wrap(errType, &err, "MGet")
	return client.client.MGet(ctx, keys...).Result()
}

func (client *singleClient) Exist(ctx context.Context, key string) (result int64, err error) {
	defer errs.Wrap(errType, &err, "Exist")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Exists(ctx, key).Result()
	}

	return client.client.Exists(ctx, key).Result()
}

func (client *singleClient) GetBytes(ctx context.Context, key string) (result []byte, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Bytes()
	} else {
		result, err = client.client.Get(ctx, key).Bytes()
	}
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

func (client *singleClient) GetInt(ctx context.Context, key string) (result int, err error) {
	defer errs.Wrap(errType, &err, "GetInt")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Int()
	} else {
		result, err = client.client.Get(ctx, key).Int()
	}
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

func (client *singleClient) Parse(ctx context.Context, key string, export any) (err error) {
	defer errs.Wrap(errType, &err, "Parse")

	if err = validateKey(key); err != nil {
		return err
	}

	var result []byte
	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Bytes()
	} else {
		result, err = client.client.Get(ctx, key).Bytes()
	}
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

func (client *singleClient) HSet(ctx context.Context, key string, value map[string]any) (err error) {
	defer errs.Wrap(errType, &err, "HSet")

	if err = validateKey(key); err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HSet(ctx, key, value).Err()
	}

	return client.client.HSet(ctx, key, value).Err()
}

func (client *singleClient) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	defer errs.Wrap(errType, &err, "HGetAll")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGetAll(ctx, key).Result()
	}

	return client.client.HGetAll(ctx, key).Result()
}

func (client *singleClient) HGet(ctx context.Context, key, field string) (result string, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Result()
	}

	return client.client.HGet(ctx, key, field).Result()
}

func (client *singleClient) HGetInt(ctx context.Context, key, field string) (result int, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Int()
	}

	return client.client.HGet(ctx, key, field).Int()
}

func (client *singleClient) HGetBool(ctx context.Context, key, field string) (result bool, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Bool()
	}

	return client.client.HGet(ctx, key, field).Bool()
}

func (client *singleClient) HExist(ctx context.Context, key, field string) (exist bool, err error) {
	defer errs.Wrap(errType, &err, "HExist")

	if err = validateKey(key); err != nil {
		return exist, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HExists(ctx, key, field).Result()
	}

	return client.client.HExists(ctx, key, field).Result()
}

func (client *singleClient) Scan(ctx context.Context, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "Scan")

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Scan(ctx, cursor, pattern, count).Result()
	}

	return client.client.Scan(ctx, cursor, pattern, count).Result()
}
