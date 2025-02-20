package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/list"
	"github.com/boostgo/lite/storage"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"time"
)

type ClientSelector func(ctx context.Context, clients []ShardClient) ShardClient

type shardClient struct {
	clients *Clients
}

// NewShard creates client implementation as shard client.
//
// Need to provide Clients object which contains multiple clients for sharding
func NewShard(clients *Clients) Client {
	return &shardClient{
		clients: clients,
	}
}

func (client *shardClient) Close() error {
	return client.clients.Close()
}

func (client *shardClient) Client(ctx context.Context) (redis.UniversalClient, error) {
	c, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	return c.Client(), nil
}

func (client *shardClient) Pipeline(ctx context.Context) (redis.Pipeliner, error) {
	raw, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	return raw.Client().Pipeline(), nil
}

func (client *shardClient) TxPipeline(ctx context.Context) (redis.Pipeliner, error) {
	raw, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	return raw.Client().TxPipeline(), nil
}

func (client *shardClient) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "Keys")

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Keys(ctx, pattern).Result()
	}

	return raw.Client().Keys(ctx, pattern).Result()
}

func (client *shardClient) Delete(ctx context.Context, keys ...string) (err error) {
	if len(keys) == 0 {
		return nil
	}

	defer errs.Wrap(errType, &err, "Delete")

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

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

	return raw.Client().Del(ctx, keys...).Err()
}

func (client *shardClient) Dump(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Dump")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Dump(ctx, key).Result()
	}

	return raw.Client().Dump(ctx, key).Result()
}

func (client *shardClient) Rename(ctx context.Context, oldKey, newKey string) (err error) {
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

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Rename(ctx, oldKey, newKey).Err()
	}

	return raw.Client().Rename(ctx, oldKey, newKey).Err()
}

func (client *shardClient) Refresh(ctx context.Context, key string, ttl time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Refresh")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Expire(ctx, key, ttl).Err()
	}

	return raw.Client().Expire(ctx, key, ttl).Err()
}

func (client *shardClient) RefreshAt(ctx context.Context, key string, at time.Time) (err error) {
	defer errs.Wrap(errType, &err, "RefreshAt")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExpireAt(ctx, key, at).Err()
	}

	return raw.Client().ExpireAt(ctx, key, at).Err()
}

func (client *shardClient) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	defer errs.Wrap(errType, &err, "TTL")

	if err = validateKey(key); err != nil {
		return ttl, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return ttl, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		ttl, err = tx.TTL(ctx, key).Result()
	} else {
		ttl, err = raw.Client().TTL(ctx, key).Result()
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

func (client *shardClient) Set(ctx context.Context, key string, value any, ttl ...time.Duration) (err error) {
	defer errs.Wrap(errType, &err, "Set")

	if err = validateKey(key); err != nil {
		return err
	}

	var expireAt time.Duration
	if len(ttl) > 0 && ttl[0] > 0 {
		expireAt = ttl[0]
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Set(ctx, key, value, expireAt).Err()
	}

	return raw.Client().Set(ctx, key, value, expireAt).Err()
}

func (client *shardClient) Get(ctx context.Context, key string) (result string, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Result()
	} else {
		result, err = raw.Client().Get(ctx, key).Result()
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

func (client *shardClient) MGet(ctx context.Context, keys []string) (result []any, err error) {
	defer errs.Wrap(errType, &err, "MGet")

	validateKeys(keys)

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.MGet(ctx, keys...).Result()
	} else {
		result, err = raw.Client().MGet(ctx, keys...).Result()
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func (client *shardClient) Exist(ctx context.Context, key string) (result int64, err error) {
	defer errs.Wrap(errType, &err, "Exist")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Exists(ctx, key).Result()
	}

	return raw.Client().Exists(ctx, key).Result()
}

func (client *shardClient) GetBytes(ctx context.Context, key string) (result []byte, err error) {
	defer errs.Wrap(errType, &err, "Get")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Bytes()
	} else {
		result, err = raw.Client().Get(ctx, key).Bytes()
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

func (client *shardClient) GetInt(ctx context.Context, key string) (result int, err error) {
	defer errs.Wrap(errType, &err, "GetInt")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Int()
	} else {
		result, err = raw.Client().Get(ctx, key).Int()
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

func (client *shardClient) Parse(ctx context.Context, key string, export any) (err error) {
	defer errs.Wrap(errType, &err, "Parse")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	var result []byte
	tx, ok := GetTx(ctx)
	if ok {
		result, err = tx.Get(ctx, key).Bytes()
	} else {
		result, err = raw.Client().Get(ctx, key).Bytes()
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

func (client *shardClient) HSet(ctx context.Context, key string, value map[string]any) (err error) {
	defer errs.Wrap(errType, &err, "HSet")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HSet(ctx, key, value).Err()
	}

	return raw.Client().HSet(ctx, key, value).Err()
}

func (client *shardClient) HGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	defer errs.Wrap(errType, &err, "HGetAll")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return nil, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGetAll(ctx, key).Result()
	}

	return raw.Client().HGetAll(ctx, key).Result()
}

func (client *shardClient) HGet(ctx context.Context, key, field string) (result string, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Result()
	}

	return raw.Client().HGet(ctx, key, field).Result()
}

func (client *shardClient) HGetInt(ctx context.Context, key, field string) (result int, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Int()
	}

	return raw.Client().HGet(ctx, key, field).Int()
}

func (client *shardClient) HGetBool(ctx context.Context, key, field string) (result bool, err error) {
	defer errs.Wrap(errType, &err, "HGet")

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HGet(ctx, key, field).Bool()
	}

	return raw.Client().HGet(ctx, key, field).Bool()
}

func (client *shardClient) HExist(ctx context.Context, key, field string) (exist bool, err error) {
	defer errs.Wrap(errType, &err, "HExist")

	if err = validateKey(key); err != nil {
		return exist, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return exist, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HExists(ctx, key, field).Result()
	}

	return raw.Client().HExists(ctx, key, field).Result()
}

func (client *shardClient) HDelete(ctx context.Context, key string, fields ...string) (err error) {
	defer errs.Wrap(errType, &err, "HDelete")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HDel(ctx, key, fields...).Err()
	}

	return raw.Client().HDel(ctx, key, fields...).Err()
}

func (client *shardClient) HScan(ctx context.Context, key string, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "HScan")

	if err = validateKey(key); err != nil {
		return keys, nextCursor, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return keys, nextCursor, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HScan(ctx, key, cursor, pattern, count).Result()
	}

	return raw.Client().HScan(ctx, key, cursor, pattern, count).Result()
}

func (client *shardClient) Scan(ctx context.Context, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "Scan")

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return keys, nextCursor, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.Scan(ctx, cursor, pattern, count).Result()
	}

	return raw.Client().Scan(ctx, cursor, pattern, count).Result()
}

func (client *shardClient) HIncrBy(ctx context.Context, key, field string, incr int64) (err error) {
	defer errs.Wrap(errType, &err, "HIncrBy")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HIncrBy(ctx, key, field, incr).Err()
	}

	return raw.Client().HIncrBy(ctx, key, field, incr).Err()
}

func (client *shardClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) (err error) {
	defer errs.Wrap(errType, &err, "HIncrByFloat")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HIncrByFloat(ctx, key, field, incr).Err()
	}

	return raw.Client().HIncrByFloat(ctx, key, field, incr).Err()
}

func (client *shardClient) HKeys(ctx context.Context, key string) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "HKeys")

	if err = validateKey(key); err != nil {
		return keys, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return keys, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HKeys(ctx, key).Result()
	}

	return raw.Client().HKeys(ctx, key).Result()
}

func (client *shardClient) HLen(ctx context.Context, key string) (length int64, err error) {
	defer errs.Wrap(errType, &err, "HLen")

	if err = validateKey(key); err != nil {
		return length, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return length, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HLen(ctx, key).Result()
	}

	return raw.Client().HLen(ctx, key).Result()
}

func (client *shardClient) HMGet(ctx context.Context, key string, fields ...string) (result []any, err error) {
	defer errs.Wrap(errType, &err, "HMGet")

	if len(fields) == 0 {
		return []any{}, nil
	}

	if err = validateKey(key); err != nil {
		return result, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return result, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HMGet(ctx, key, fields...).Result()
	}

	return raw.Client().HMGet(ctx, key, fields...).Result()
}

func (client *shardClient) HMSet(ctx context.Context, key string, values ...any) (err error) {
	defer errs.Wrap(errType, &err, "HMSet")

	if len(values) == 0 {
		return nil
	}

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HMSet(ctx, key, values...).Err()
	}

	return raw.Client().HMSet(ctx, key, values...).Err()
}

func (client *shardClient) HSetNX(ctx context.Context, key, field string, value any) (err error) {
	defer errs.Wrap(errType, &err, "HSetNX")

	if err = validateKey(key); err != nil {
		return err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HSetNX(ctx, key, field, value).Err()
	}

	return raw.Client().HSetNX(ctx, key, field, value).Err()
}

func (client *shardClient) HScanNoValues(ctx context.Context, key string, cursor uint64, pattern string, count int64) (keys []string, nextCursor uint64, err error) {
	defer errs.Wrap(errType, &err, "HScaNoValues")

	if err = validateKey(key); err != nil {
		return keys, cursor, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return keys, nextCursor, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HScanNoValues(ctx, key, cursor, pattern, count).Result()
	}

	return raw.Client().HScanNoValues(ctx, key, cursor, pattern, count).Result()
}

func (client *shardClient) HVals(ctx context.Context, key string) (values []string, err error) {
	defer errs.Wrap(errType, &err, "HVals")

	if err = validateKey(key); err != nil {
		return values, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return values, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HVals(ctx, key).Result()
	}

	return raw.Client().HVals(ctx, key).Result()
}

func (client *shardClient) HRandField(ctx context.Context, key string, count int) (keys []string, err error) {
	defer errs.Wrap(errType, &err, "HRandField")

	if err = validateKey(key); err != nil {
		return keys, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return keys, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HRandField(ctx, key, count).Result()
	}

	return raw.Client().HRandField(ctx, key, count).Result()
}

func (client *shardClient) HRandFieldWithValues(ctx context.Context, key string, count int) (values []redis.KeyValue, err error) {
	defer errs.Wrap(errType, &err, "HRandFieldWithValues")

	if err = validateKey(key); err != nil {
		return values, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return values, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HRandFieldWithValues(ctx, key, count).Result()
	}

	return raw.Client().HRandFieldWithValues(ctx, key, count).Result()
}

func (client *shardClient) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) (states []int64, err error) {
	defer errs.Wrap(errType, &err, "HExpire")

	if len(fields) == 0 {
		return states, nil
	}

	if err = validateKey(key); err != nil {
		return states, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return states, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HExpire(ctx, key, expiration, fields...).Result()
	}

	return raw.Client().HExpire(ctx, key, expiration, fields...).Result()
}

func (client *shardClient) HTTL(ctx context.Context, key string, fields ...string) (ttls []int64, err error) {
	defer errs.Wrap(errType, &err, "HTTL")

	if len(fields) == 0 {
		return ttls, nil
	}

	if err = validateKey(key); err != nil {
		return ttls, err
	}

	raw, err := client.clients.Get(ctx)
	if err != nil {
		return ttls, err
	}

	tx, ok := GetTx(ctx)
	if ok {
		return tx.HTTL(ctx, key, fields...).Result()
	}

	return raw.Client().HTTL(ctx, key, fields...).Result()
}

type ShardClient interface {
	Key() string
	Conditions() []string
	Client() redis.UniversalClient
	Close() error
}

// Clients contain all clients for shard client and selector for choosing connection
type Clients struct {
	clients  []ShardClient
	selector ClientSelector
}

func newClients(clients []ShardClient, selector ClientSelector) *Clients {
	return &Clients{
		clients:  clients,
		selector: selector,
	}
}

// Get returns shard connect by using selector
func (c *Clients) Get(ctx context.Context) (ShardClient, error) {
	// get shard by provided selector
	conn := c.selector(ctx, c.clients)
	if conn == nil {
		return nil, storage.ErrConnNotSelected
	}

	return conn, nil
}

// Clients return all shard clients
func (c *Clients) Clients() []ShardClient {
	return c.clients
}

// RawConnections returns all clients as []*sqlx.DB
func (c *Clients) RawConnections() []redis.UniversalClient {
	clients := make([]redis.UniversalClient, len(c.clients))
	for idx, client := range c.clients {
		clients[idx] = client.Client()
	}
	return clients
}

// Close all clients in parallel
func (c *Clients) Close() error {
	wg := errgroup.Group{}

	for _, conn := range c.clients {
		wg.Go(conn.Close)
	}

	return wg.Wait()
}
