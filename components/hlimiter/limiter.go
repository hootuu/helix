package hlimiter

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/hootuu/helix/components/zplt"
)

type Result = redis_rate.Result

type Limiter struct {
	limiter *redis_rate.Limiter
}

func NewLimiter() *Limiter {
	hrds := zplt.HelixRdsCache()
	limiter := redis_rate.NewLimiter(hrds.Redis())
	return &Limiter{limiter: limiter}
}

func (r *Limiter) Limiter() *redis_rate.Limiter {
	return r.limiter
}

func (r *Limiter) Allow(key string, limitInSec int) bool {
	ctx := context.Background()
	result, err := r.limiter.Allow(ctx, key, redis_rate.PerSecond(limitInSec))
	if err != nil {
		return false
	}
	if result.Allowed > 0 {
		return true
	}
	return false
}
