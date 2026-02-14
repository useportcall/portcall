package ratelimitx

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Check checks if a request is allowed based on the rate limit.
// key: unique identifier (e.g., "api:userid123", "ip:192.168.1.1")
func (rl *RateLimiter) Check(key string, config Config) (*Result, error) {
	if err := rl.client.Ping(rl.ctx).Err(); err != nil {
		log.Printf("Redis unavailable, allowing request: %v", err)
		return allowResult(config), nil
	}

	now := time.Now()
	windowStart := now.Truncate(config.Window)
	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, windowStart.Unix())

	count, err := rl.runIncrScript(redisKey, int(config.Window.Seconds()))
	if err != nil {
		log.Printf("Redis script execution failed: %v", err)
		return allowResult(config), nil
	}

	allowed := count <= int64(config.Limit)
	remaining := config.Limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return &Result{
		Allowed:   allowed,
		Limit:     config.Limit,
		Remaining: remaining,
		ResetAt:   windowStart.Add(config.Window),
	}, nil
}

func (rl *RateLimiter) runIncrScript(key string, ttlSecs int) (int64, error) {
	script := redis.NewScript(`
		local current = redis.call("INCR", KEYS[1])
		if current == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[1])
		end
		return current
	`)
	result, err := script.Run(rl.ctx, rl.client, []string{key}, ttlSecs).Result()
	if err != nil {
		return 0, err
	}
	count, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected Redis result type: %T", result)
	}
	return count, nil
}
