package hrds

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/hootuu/helix/helix"
	"github.com/hootuu/hyle/hcfg"
	"github.com/hootuu/hyle/hlog"
	"github.com/hootuu/hyle/hsys"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type Cache struct {
	code   string
	client *redis.Client
}

func newCache(code string) *Cache {
	return &Cache{
		code:   code,
		client: nil,
	}
}

func (cache *Cache) Helix() helix.Helix {
	return helix.BuildHelix(cache.code, cache.startup, cache.shutdown)
}

func (cache *Cache) Redis() *redis.Client {
	return cache.client
}

func (cache *Cache) Check() error {
	err := retry.Do(
		func() error {
			statusCmd := cache.client.Ping(context.Background())
			if statusCmd.Err() != nil {
				hsys.Error("# Connecting to redis [", cache.code, "] error:"+statusCmd.Err().Error()+"#\n")
				return statusCmd.Err()
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(3*time.Second),
	)
	if err != nil {
		return err
	}
	return nil
}

func (cache *Cache) startup() (context.Context, error) {
	hsys.Info("\n# Connecting to redis [", cache.code, "] ... #")
	cache.client = redis.NewClient(&redis.Options{
		Network:                    hcfg.GetString("redis."+cache.code+".network", ""),
		Addr:                       hcfg.GetString("redis."+cache.code+".addr", "127.0.0.1:6379"),
		ClientName:                 hcfg.GetString("redis."+cache.code+".client.name", ""),
		Dialer:                     nil,
		OnConnect:                  nil,
		Protocol:                   hcfg.GetInt("redis."+cache.code+".protocol", 0),
		Username:                   hcfg.GetString("redis."+cache.code+".username", ""),
		Password:                   hcfg.GetString("redis."+cache.code+".password", ""),
		CredentialsProvider:        nil,
		CredentialsProviderContext: nil,
		DB:                         hcfg.GetInt("redis."+cache.code+".db", 0),
		MaxRetries:                 hcfg.GetInt("redis."+cache.code+".max.retries", 0),
		MinRetryBackoff:            hcfg.GetDuration("redis."+cache.code+".min.retry.backoff", 0),
		MaxRetryBackoff:            hcfg.GetDuration("redis."+cache.code+".max.retry.backoff", 0),
		DialTimeout:                hcfg.GetDuration("redis."+cache.code+".dial.timeout", 0),
		ReadTimeout:                hcfg.GetDuration("redis."+cache.code+".read.timeout", 0),
		WriteTimeout:               hcfg.GetDuration("redis."+cache.code+".write.timeout", 0),
		ContextTimeoutEnabled:      hcfg.GetBool("redis."+cache.code+".context.timeout.enabled", false),
		PoolFIFO:                   hcfg.GetBool("redis."+cache.code+".pool.fifo", false),
		PoolSize:                   hcfg.GetInt("redis."+cache.code+".pool.size", 0),
		PoolTimeout:                hcfg.GetDuration("redis."+cache.code+".pool.timeout", 0),
		MinIdleConns:               hcfg.GetInt("redis."+cache.code+".min.idle.conns", 0),
		MaxIdleConns:               hcfg.GetInt("redis."+cache.code+".max.idle.conns", 0),
		MaxActiveConns:             hcfg.GetInt("redis."+cache.code+".max.active.conns", 0),
		ConnMaxIdleTime:            hcfg.GetDuration("redis."+cache.code+".conn.max.idle.time", 0),
		ConnMaxLifetime:            hcfg.GetDuration("redis."+cache.code+".conn.max.life.time", 0),
		TLSConfig:                  nil,
		Limiter:                    nil,
		DisableIdentity:            hcfg.GetBool("redis."+cache.code+".disable.identity", false),
		IdentitySuffix:             hcfg.GetString("redis."+cache.code+".identity.suffix", ""),
		UnstableResp3:              hcfg.GetBool("redis."+cache.code+".unstable.resp3", false),
	})
	hsys.Success("# Connecting to redis [", cache.code, "] OK #\n")
	return context.Background(), nil
}

func (cache *Cache) shutdown(_ context.Context) {
	if cache.client != nil {
		_ = cache.client.Close()
	}
}

var gRedisMap = make(map[string]*Cache)
var gRedisMutex sync.Mutex

func doRegister(code string) {
	gRedisMutex.Lock()
	defer gRedisMutex.Unlock()
	if _, ok := gRedisMap[code]; ok {
		hlog.Err("hrds.doRegister: redis repetition")
		return
	}
	cache := newCache(code)
	gRedisMap[code] = cache
	helix.Use(cache.Helix())
}

func doGetCache(code string) *Cache {
	gRedisMutex.Lock()
	defer gRedisMutex.Unlock()
	cache, ok := gRedisMap[code]
	if !ok {
		return nil
	}
	return cache
}

func doCacheActWithRetry(call func() error) error {
	return retry.Do(
		call,
		retry.Attempts(uint(hcfg.GetInt("hrds.act.retry.attempts", 3))),
		retry.Delay(hcfg.GetDuration("hrds.act.retry.delay", 100*time.Millisecond)),
	)
}
