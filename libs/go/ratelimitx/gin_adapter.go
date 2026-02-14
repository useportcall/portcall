package ratelimitx

import (
	"fmt"
	"strconv"
	"time"
)

// GinAdapter provides a Gin-compatible middleware wrapper for rate limiting
type GinAdapter struct {
	limiter *RateLimiter
}

// NewGinAdapter creates a new Gin middleware adapter
func NewGinAdapter() *GinAdapter {
	return &GinAdapter{limiter: New()}
}

// GinContext is a minimal interface for Gin context.
// This allows us to work with gin.Context without importing it.
type GinContext interface {
	Get(key any) (value interface{}, exists bool)
	GetString(key any) string
	GetHeader(key string) string
	ClientIP() string
	Next()
	AbortWithStatusJSON(code int, jsonObj interface{})
	Header(key, value string)
}

// Middleware returns a Gin-compatible middleware function.
func (ga *GinAdapter) Middleware(keyExtractor func(GinContext) string, limit int, window time.Duration) func(GinContext) {
	config := Config{Limit: limit, Window: window}

	return func(c GinContext) {
		key := keyExtractor(c)
		if key == "" {
			c.Next()
			return
		}

		result, err := ga.limiter.Check(key, config)
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(result.Limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

		if !result.Allowed {
			c.Header("Retry-After", strconv.FormatInt(int64(time.Until(result.ResetAt).Seconds()), 10))
			c.AbortWithStatusJSON(429, map[string]interface{}{
				"error": fmt.Sprintf("Rate limit exceeded. Try again at %s", result.ResetAt.Format("2006-01-02T15:04:05Z07:00")),
			})
			return
		}

		c.Next()
	}
}
