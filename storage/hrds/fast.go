package hrds

import (
	"fmt"
	"time"
)

func Fast[T any](
	cache *Cache,
	code string,
	key string,
	duration time.Duration,
	load func() (*T, error)) (*T, error) {
	uniKey := fmt.Sprintf("%s:%s", code, key)
	m, err := CacheGet[T](cache, uniKey)
	if err != nil {
		return nil, err
	}
	if m == nil {
		m, err = load()
		if err != nil {
			return nil, err
		}
		if m == nil {
			return nil, nil
		}
		// TODO add
		//err = cache.Set(uniKey, m, duration)
		//if err != nil {
		//	hlog.Fix("hrds.Fast.cache.Set failed", zap.Error(err))
		//}
		return m, nil
	}
	return m, nil
}
