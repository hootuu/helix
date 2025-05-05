package hidem

import (
	"context"
	"fmt"
	"github.com/hootuu/helix/components/zplt"
	"github.com/hootuu/hyle/hlog"
	"go.uber.org/zap"
	"time"
)

type cacheFactory struct {
	code       string
	expiration time.Duration
}

func (f *cacheFactory) cacheKey(idemCode string) string {
	return fmt.Sprintf("%s:%s", f.code, idemCode)
}

func (f *cacheFactory) Check(idemCode string) (bool, error) {
	if err := CheckIdemCode(idemCode); err != nil {
		return false, err
	}
	ok, err := zplt.HelixRdsCache().Redis().SetNX(
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

func newCacheFactory(code string, expiration time.Duration) (*cacheFactory, error) {
	f := &cacheFactory{
		code:       code,
		expiration: expiration,
	}
	return f, nil
}
