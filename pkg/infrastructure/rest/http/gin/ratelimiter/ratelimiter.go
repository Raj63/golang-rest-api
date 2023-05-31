package ratelimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// GinRatelimiterConfig Gin Ratelimiter
type GinRatelimiterConfig struct {
	// LimitKey
	LimitKey func(*gin.Context) string
	// LimitedHandler
	LimitedHandler func(*gin.Context)
	// TokenBucketConfig
	TokenBucketConfig func(*gin.Context) (time.Duration, int)
}

// DefaultGinLimitKey returns the default gin limit key
func DefaultGinLimitKey(c *gin.Context) string {
	return fmt.Sprintf("pink-lady:ratelimiter:%s:%s", c.ClientIP(), c.FullPath())
}

// DefaultGinLimitedHandler 429
func DefaultGinLimitedHandler(c *gin.Context) {
	c.AbortWithStatus(http.StatusTooManyRequests)
}

// GinMemRatelimiter return the gin rate limiter handler func
func GinMemRatelimiter(conf GinRatelimiterConfig, logger *logger.Logger) gin.HandlerFunc {
	if conf.TokenBucketConfig == nil {
		panic("GinRatelimiterConfig must implement the TokenBucketConfig callback function")
	}
	limiter := NewMemRatelimiter(logger)

	return func(c *gin.Context) {
		// limit key
		limitKey := DefaultGinLimitKey(c)
		if conf.LimitKey != nil {
			limitKey = conf.LimitKey(c)
		}

		limitedHandler := DefaultGinLimitedHandler
		if conf.LimitedHandler != nil {
			limitedHandler = conf.LimitedHandler
		}

		tokenFillInterval, bucketSize := conf.TokenBucketConfig(c)

		if !limiter.Allow(c, limitKey, tokenFillInterval, bucketSize) {
			limitedHandler(c)
			return
		}
		c.Next()
	}
}
