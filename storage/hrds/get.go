package hrds

import (
	"context"
	"errors"
	"github.com/hootuu/hyle/data/hjson"
	"github.com/redis/go-redis/v9"
)

func (cache *Cache) Get(key string, parse func(cmd *redis.StringCmd)) (bool, error) {
	r := cache.Redis().Get(context.Background(), key)
	if err := r.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	parse(r)
	return true, nil
}

func (cache *Cache) GetString(key string) (string, error) {
	str := ""
	_, _ = cache.Get(key, func(cmd *redis.StringCmd) {
		str = cmd.String()
	})
	return str, nil
}

func (cache *Cache) GetInt64(key string) (int64, error) {
	val := int64(0)
	var err error
	_, _ = cache.Get(key, func(cmd *redis.StringCmd) {
		val, err = cmd.Int64()
	})
	return val, err
}

func CacheGet[T any](cache *Cache, key string) (*T, error) {
	var obj *T
	var err error
	_, _ = cache.Get(key, func(cmd *redis.StringCmd) {
		var bytes []byte
		bytes, err = cmd.Bytes()
		if err != nil {
			return
		}
		if len(bytes) == 0 {
			return
		}
		obj, err = hjson.FromBytes[T](bytes)
	})
	return obj, err
}
