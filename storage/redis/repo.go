package redis

import (
	"context"
	"encoding/json"
	"github.com/boostgo/lite/errs"
	"github.com/redis/go-redis/v9"
	"time"
)

const errType = "Redis"

type Repository struct {
	conn *redis.Client
}

func New(address string, port, db int, password string, opts ...Option) (*Repository, error) {
	conn, err := Connect(address, port, db, password, opts...)
	if err != nil {
		return nil, err
	}

	return &Repository{
		conn: conn,
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
		conn: conn,
	}
}

func (repo Repository) Close() error {
	return repo.conn.Close()
}

func (repo Repository) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "Keys")
	return repo.conn.Keys(ctx, pattern).Result()
}

func (repo Repository) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return nil
	}

	defer errs.Wrap(errType, &err, "Delete")
	return repo.conn.Del(ctx, keys...).Err()
}

func (repo Repository) Rename(ctx context.Context, oldKey, newKey string) (err error) {
	defer errs.Wrap(errType, &err, "Rename")
	return repo.conn.Rename(ctx, oldKey, newKey).Err()
}

func (repo Repository) Refresh(ctx context.Context, key string, ttl time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Refresh")
	return repo.conn.Expire(ctx, key, ttl).Err()
}

func (repo Repository) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	defer errs.Wrap(errType, &err, "TTL")
	return repo.conn.TTL(ctx, key).Result()
}

func (repo Repository) Set(ctx context.Context, key string, value any, ttl ...time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Set")

	var expireAt time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = ttl[0]
	}

	return repo.conn.Set(ctx, key, value, expireAt).Err()
}

func (repo Repository) Get(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Get")
	return repo.conn.Get(ctx, key).Result()
}

func (repo Repository) GetBytes(ctx context.Context, key string) (result []byte, err error) {
	defer errs.Wrap(errType, &err, "Get")
	return repo.conn.Get(ctx, key).Bytes()
}

func (repo Repository) GetInt(ctx context.Context, key string) (result int, err error) {
	defer errs.Wrap(errType, &err, "GetInt")
	return repo.conn.Get(ctx, key).Int()
}

func (repo Repository) Parse(ctx context.Context, key string, export any) (err error) {
	defer errs.Wrap(errType, &err, "Parse")

	var result []byte
	result, err = repo.conn.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(result, &export)
}

func (repo Repository) HSet(ctx context.Context, key string, value map[string]any) (err error) {
	defer errs.Wrap(errType, &err, "HSet")
	return repo.conn.HSet(ctx, key, value).Err()
}

func (repo Repository) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	defer errs.Wrap(errType, &err, "HGetAll")
	return repo.conn.HGetAll(ctx, key).Result()
}

func (repo Repository) HGet(ctx context.Context, key, field string) (result string, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.conn.HGet(ctx, key, field).Result()
}

func (repo Repository) HGetInt(ctx context.Context, key, field string) (result int, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.conn.HGet(ctx, key, field).Int()
}

func (repo Repository) HGetBool(ctx context.Context, key, field string) (result bool, err error) {
	defer errs.Wrap(errType, &err, "HGet")
	return repo.conn.HGet(ctx, key, field).Bool()
}
