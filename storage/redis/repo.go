package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/redis/go-redis/v9"
	"time"
)

const errType = "Redis"

type Repository struct {
	client *redis.Client
}

func New(address string, port, db int, password string, opts ...Option) (*Repository, error) {
	conn, err := Connect(address, port, db, password, opts...)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: conn,
	}, nil
}

func Must(address string, port, db int, password string, opts ...Option) *Repository {
	repo, err := New(address, port, db, password, opts...)
	if err != nil {
		panic(err)
	}

	return repo
}

func NewFromClient(conn *redis.Client) *Repository {
	return &Repository{
		client: conn,
	}
}

func (repo Repository) Close() error {
	return repo.client.Close()
}

func (repo Repository) Client() *redis.Client {
	return repo.client
}

func (repo Repository) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "Keys")
	return repo.client.Keys(ctx, pattern).Result()
}

func (repo Repository) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return nil
	}

	defer errs.Wrap(errType, &err, "Delete")
	return repo.client.Del(ctx, keys...).Err()
}

func (repo Repository) Dump(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Dump")
	return repo.client.Dump(ctx, key).Result()
}

func (repo Repository) Rename(ctx context.Context, oldKey, newKey string) (err error) {
	defer errs.Wrap(errType, &err, "Rename")
	return repo.client.Rename(ctx, oldKey, newKey).Err()
}

func (repo Repository) Refresh(ctx context.Context, key string, ttl time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Refresh")
	return repo.client.Expire(ctx, key, ttl).Err()
}

func (repo Repository) RefreshAt(ctx context.Context, key string, at time.Time) (err error) {
	defer errs.Wrap(errType, &err, "RefreshAt")
	return repo.client.ExpireAt(ctx, key, at).Err()
}

func (repo Repository) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	defer errs.Wrap(errType, &err, "TTL")

	ttl, err = repo.client.TTL(ctx, key).Result()
	if err != nil {
		return ttl, err
	}

	const notExistKey = -2
	if ttl == notExistKey {
		return ttl, errs.ErrNotFound
	}

	return ttl, nil
}

func (repo Repository) Set(ctx context.Context, key string, value any, ttl ...time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Set")

	var expireAt time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = ttl[0]
	}

	return repo.client.Set(ctx, key, value, expireAt).Err()
}

func (repo Repository) Get(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Get")

	result, err = repo.client.Get(ctx, key).Result()
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

func (repo Repository) Exist(ctx context.Context, key string) (result int64, err error) {
	defer errs.Wrap(errType, &err, "Exist")
	return repo.client.Exists(ctx, key).Result()
}

func (repo Repository) GetBytes(ctx context.Context, key string) (result []byte, err error) {
	defer errs.Wrap(errType, &err, "Get")

	result, err = repo.client.Get(ctx, key).Bytes()
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

func (repo Repository) GetInt(ctx context.Context, key string) (result int, err error) {
	defer errs.Wrap(errType, &err, "GetInt")

	result, err = repo.client.Get(ctx, key).Int()
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

func (repo Repository) Parse(ctx context.Context, key string, export any) (err error) {
	defer errs.Wrap(errType, &err, "Parse")

	var result []byte
	result, err = repo.client.Get(ctx, key).Bytes()
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

func (repo Repository) HSet(ctx context.Context, key string, value map[string]any) (err error) {
	defer errs.Wrap(errType, &err, "HSet")
	return repo.client.HSet(ctx, key, value).Err()
}

func (repo Repository) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	defer errs.Wrap(errType, &err, "HGetAll")
	return repo.client.HGetAll(ctx, key).Result()
}

func (repo Repository) HGet(ctx context.Context, key, field string) (result string, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.client.HGet(ctx, key, field).Result()
}

func (repo Repository) HGetInt(ctx context.Context, key, field string) (result int, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.client.HGet(ctx, key, field).Int()
}

func (repo Repository) HGetBool(ctx context.Context, key, field string) (result bool, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.client.HGet(ctx, key, field).Bool()
}

func (repo Repository) HExist(ctx context.Context, key, field string) (exist bool, err error) {
	defer errs.Wrap(errType, &err, "HExist")
	return repo.client.HExists(ctx, key, field).Result()
}

func (repo Repository) Scan(ctx context.Context, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "Scan")
	return repo.client.Scan(ctx, cursor, pattern, count).Result()
}
