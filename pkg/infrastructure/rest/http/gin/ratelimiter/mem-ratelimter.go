package ratelimiter

import (
	"context"
	"time"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

// MemRatelimiter limiter
type MemRatelimiter struct {
	*rate.Limiter
	*cache.Cache
	Expire time.Duration
	logger *logger.Logger
}

var (
	// MemRatelimiterCacheExpiration MemRatelimiter key
	MemRatelimiterCacheExpiration = time.Minute * 60
	// MemRatelimiterCacheCleanInterval MemRatelimiter
	MemRatelimiterCacheCleanInterval = time.Minute * 60
)

// NewMemRatelimiter mem limiter
func NewMemRatelimiter(logger *logger.Logger) *MemRatelimiter {
	// mem cache
	memCache := cache.New(MemRatelimiterCacheExpiration, MemRatelimiterCacheCleanInterval)
	return &MemRatelimiter{
		Cache:  memCache,
		logger: logger,
	}
}

// Allow time/rate & token bucket key
// tokenFillInterval Token
// bucketSize Token
func (r *MemRatelimiter) Allow(ctx context.Context, key string, tokenFillInterval time.Duration, bucketSize int) bool {
	if tokenFillInterval.Seconds() <= 0 || bucketSize <= 0 {
		return false
	}

	tokenRate := rate.Every(tokenFillInterval)
	limiterI, exists := r.Cache.Get(key)
	if !exists {
		limiter := rate.NewLimiter(tokenRate, bucketSize)
		limiter.Allow()
		r.Cache.Set(key, limiter, MemRatelimiterCacheExpiration)
		return true
	}

	if limiter, ok := limiterI.(*rate.Limiter); ok {
		isAllow := limiter.Allow()
		r.Cache.Set(key, limiter, MemRatelimiterCacheExpiration)
		return isAllow
	}

	r.logger.Errorf("MemRatelimiter assert limiter error")
	return true

}
