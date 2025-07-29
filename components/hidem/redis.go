package hidem

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/storage/hrds"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"time"
)

type cacheFactory struct {
	code       string
	expiration time.Duration
	cache      *hrds.Cache
}

func (f *cacheFactory) cacheKey(idemCode string) string {
	return fmt.Sprintf("%s:%s", f.code, idemCode)
}

func (f *cacheFactory) Check(idemCode string) (bool, error) {
	if err := CheckIdemCode(idemCode); err != nil {
		return false, err
	}
	ok, err := f.cache.Redis().SetNX(
		context.Background(),
		idemCode,
		true,
		f.expiration,
	).Result()
	if err != nil {
		hlog.Err("helix.idem.redis.Check", zap.Error(err))
		return false, err
	}
	return ok, nil
}

func (f *cacheFactory) MustCheck(idemCode string) error {
	ok, err := f.Check(idemCode)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("idem error: %s", idemCode)
	}
	return nil
}

func newCacheFactory(cache *hrds.Cache, code string, expiration time.Duration) (*cacheFactory, error) {
	f := &cacheFactory{
		code:       code,
		expiration: expiration,
		cache:      cache,
	}
	return f, nil
}
