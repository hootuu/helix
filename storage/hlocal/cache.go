package hlocal

import (
	"github.com/hootuu/hyle/hlog"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	NoExpiration      = cache.NoExpiration
	DefaultExpiration = cache.DefaultExpiration
)

type Cache[T any] struct {
	core *cache.Cache
}

func NewCache[T any](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	return &Cache[T]{
		core: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (cache *Cache[T]) Set(key string, obj interface{}) {
	cache.core.Set(key, obj, DefaultExpiration)
}

func (cache *Cache[T]) Get(key string) *T {
	obj, ok := cache.core.Get(key)
	if !ok {
		return nil
	}
	tObj, ok := obj.(*T)
	if !ok {
		hlog.Err("hlocal.Cache.Get: invalid object type")
		return nil
	}
	return tObj
}

func (cache *Cache[T]) GetSet(key string, set func() (*T, error)) (*T, error) {
	tObj := cache.Get(key)
	if tObj == nil {
		newObj, err := set()
		if err != nil {
			return nil, err
		}
		cache.Set(key, newObj)
		return newObj, nil

	}
	return tObj, nil
}
