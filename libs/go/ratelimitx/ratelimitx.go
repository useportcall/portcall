package ratelimitx

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiter manages rate limiting using a token bucket algorithm with Redis
type RateLimiter struct {
	client *redis.Client
	ctx    context.Context
}

// Config holds the rate limiting configuration
type Config struct {
	// Limit is the maximum number of requests allowed
	Limit int
	// Window is the time window for the limit (e.g., 1 hour, 1 minute)
	Window time.Duration
}

// Result contains the rate limiting check result
type Result struct {
	Allowed   bool
	Limit     int
	Remaining int
	ResetAt   time.Time
}

// New creates a new rate limiter instance.
// Reads REDIS_ADDR, REDIS_PASSWORD, and REDIS_TLS environment variables.
func New() *RateLimiter {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTLS := os.Getenv("REDIS_TLS")

	opts := &redis.Options{Addr: redisAddr}
	if redisPassword != "" {
		opts.Password = redisPassword
	}
	if redisTLS == "true" {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis connection failed: %v. Rate limiting will be disabled.", err)
	}

	return &RateLimiter{client: client, ctx: ctx}
}

// Close closes the Redis connection.
func (rl *RateLimiter) Close() error { return rl.client.Close() }

func allowResult(config Config) *Result {
	return &Result{
		Allowed:   true,
		Limit:     config.Limit,
		Remaining: config.Limit,
		ResetAt:   time.Now().Add(config.Window),
	}
}
