package redis

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

const lockerErrType = "Redis Locker"

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

type Locker struct {
	client *redsync.Redsync
}

func NewLocker(client Client) (*Locker, error) {
	universalClient, err := client.Client(context.Background())
	if err != nil {
		return nil, err
	}

	return &Locker{
		client: redsync.New(goredis.NewPool(universalClient)),
	}, nil
}

func MustLocker(client Client) *Locker {
	locker, err := NewLocker(client)
	if err != nil {
		panic(err)
	}

	return locker
}

func (locker *Locker) lock(ctx context.Context, lockKey string) (mx Mutex, err error) {
	defer errs.Wrap(lockerErrType, &err, "lock")
	
	redisMx := locker.client.NewMutex(lockKey)
	mx = newLockMutex(redisMx)
	if err = mx.Lock(ctx); err != nil {
		return nil, err
	}

	return mx, nil
}

type lockMutex struct {
	mx *redsync.Mutex
}

func newLockMutex(mx *redsync.Mutex) Mutex {
	return &lockMutex{
		mx: mx,
	}
}

func (mx *lockMutex) Unlock(ctx context.Context) (err error) {
	_, err = mx.mx.UnlockContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (mx *lockMutex) Lock(ctx context.Context) (err error) {
	return mx.mx.LockContext(ctx)
}
