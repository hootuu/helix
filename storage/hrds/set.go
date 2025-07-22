package hrds

import (
	"context"
	"errors"
	"github.com/avast/retry-go"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

func (cache *Cache) SetString(key string, val string, expiration time.Duration) error {
	err := retry.Do(func() error {
		cmd := cache.Redis().Set(context.Background(), key, val, expiration)
		if err := cmd.Err(); err != nil {
			return err
		}
		return nil
	},
		retry.Attempts(cast.ToUint(hcfg.GetInt("hrds."+cache.code+".retry.attempts", 3))),
		retry.Delay(hcfg.GetDuration("hrds."+cache.code+".retry.delay", 200*time.Millisecond)),
	)
	if err != nil {
		hlog.Err("hrds.Set", zap.Error(err))
		return err
	}
	return nil
}

func (cache *Cache) Set(key string, val interface{}, expiration time.Duration) error {
	if val == nil {
		return errors.New("val is nil")
	}
	payload := hjson.MustToBytes(val)
	if payload == nil {
		return errors.New("payload is nil")
	}
	err := retry.Do(func() error {
		cmd := cache.Redis().Set(context.Background(), key, payload, expiration)
		if err := cmd.Err(); err != nil {
			return err
		}
		return nil
	},
		retry.Attempts(cast.ToUint(hcfg.GetInt("hrds."+cache.code+".retry.attempts", 3))),
		retry.Delay(hcfg.GetDuration("hrds."+cache.code+".retry.delay", 200*time.Millisecond)),
	)
	if err != nil {
		hlog.Err("hrds.Set", zap.Error(err))
		return err
	}
	return nil
}
