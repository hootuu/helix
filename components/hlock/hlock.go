package hlock

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type Locker struct {
	cache *hrds.Cache
}

func NewLocker(cache *hrds.Cache) *Locker {
	return &Locker{cache: cache}
}

func Light() *Locker {
	return NewLocker(zplt.HelixRdsCache())
}

func (l *Locker) Lock(
	ctx context.Context,
	key string,
	call func() error,
	ttl time.Duration,
) (bool, error) {
	rds := l.cache.Redis()

	lockKey := fmt.Sprintf("HLOCK:LOCK:%s", key)

	token := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))

	locked, err := rds.SetNX(ctx, lockKey, token, ttl).Result()
	if err != nil {
		hlog.Err("hlock.Lock: Set Lock Token", zap.String("key", key), zap.Error(err))
		return false, err
	}

	if !locked {
		hlog.Info("hlock.Lock: !locked", zap.String("key", key))
		return false, nil
	}

	defer func() {
		l.release(ctx, lockKey, token)
	}()

	err = call()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *Locker) LockWait(
	ctx context.Context,
	key string,
	call func() error,
	ttl time.Duration,
	waitTimeout time.Duration,
	retryInterval time.Duration,
) (bool, error) {
	rds := l.cache.Redis()

	lockKey := fmt.Sprintf("HLOCK:LOCK:%s", key)
	token := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))

	deadline := time.Now().Add(waitTimeout)

	for {
		locked, err := rds.SetNX(ctx, lockKey, token, ttl).Result()
		if err != nil {
			hlog.Err("hlock.LockWait: SetNX error",
				zap.String("key", key),
				zap.Error(err),
			)
			return false, err
		}

		if locked {
			defer func() {
				l.release(ctx, lockKey, token)
			}()

			err = call()
			if err != nil {
				return false, err
			}
			return true, nil
		}

		if time.Now().After(deadline) {
			hlog.Info("hlock.LockWait: timeout waiting for lock",
				zap.String("key", key),
			)
			return false, fmt.Errorf("lock wait timeout for key=%s", key)
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryInterval + time.Duration(rand.Intn(50))*time.Millisecond):
			continue
		}
	}
}

func (l *Locker) OnceLock(
	ctx context.Context,
	key string,
	call func() error,
	ttl time.Duration,
) (bool, error) {
	rds := l.cache.Redis()

	taskStatusKey := fmt.Sprintf("HLOCK:ONCE:STATUS:%s", key)

	if isDone, _ := rds.Get(ctx, taskStatusKey).Result(); isDone == "1" {
		return true, nil
	}

	lockKey := fmt.Sprintf("HLOCK:ONCE:LOCK:%s", key)

	token := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))

	locked, err := rds.SetNX(ctx, lockKey, token, ttl).Result()
	if err != nil {
		hlog.Err("hlock.Lock: Set Lock Token", zap.String("key", key), zap.Error(err))
		return false, err
	}

	if !locked {
		hlog.Info("hlock.Lock: !locked", zap.String("key", key))
		return false, nil
	}

	defer func() {
		l.release(ctx, lockKey, token)
	}()

	if isDone, _ := rds.Get(ctx, taskStatusKey).Result(); isDone == "1" {
		return true, nil
	}

	err = call()
	if err != nil {
		return false, err
	}

	if err := rds.Set(ctx, taskStatusKey, "1", ttl).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (l *Locker) release(ctx context.Context, lockKey string, token string) {
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`
	_, err := l.cache.Redis().Eval(ctx, script, []string{lockKey}, token).Result()
	if err != nil {
		hlog.Err("hlock.Lock: failed to release lock:",
			zap.String("lockKey", lockKey),
			zap.String("token", token),
			zap.Error(err))
	}
}
