package ratelimitx

import (
	"strconv"
	"strings"
)

// ByAPIKeyHeader extracts the API key from a header for rate limiting.
func ByAPIKeyHeader(c GinContext) string {
	apiKey := c.GetHeader("x-api-key")
	if apiKey == "" {
		return ""
	}
	parts := strings.Split(apiKey, "_")
	if len(parts) >= 2 {
		return "apikey:" + strings.Join(parts[0:2], "_")
	}
	return "apikey:" + apiKey
}

// ByClientIP extracts the client IP.
func ByClientIP(c GinContext) string {
	ip := c.ClientIP()
	if ip == "" {
		return ""
	}
	return "ip:" + ip
}

// ByContextKey creates a key extractor that reads from Gin context.
func ByContextKey(contextKey string, prefix string) func(GinContext) string {
	return func(c GinContext) string {
		value, exists := c.Get(contextKey)
		if !exists {
			return ""
		}
		valueStr, ok := value.(string)
		if !ok {
			if uintVal, ok := value.(uint); ok {
				valueStr = strconv.FormatUint(uint64(uintVal), 10)
			} else {
				return ""
			}
		}
		if valueStr == "" {
			return ""
		}
		return prefix + ":" + valueStr
	}
}

// ByEmail extracts email from auth context.
func ByEmail(c GinContext) string {
	email := c.GetString("auth_email")
	if email == "" {
		return ""
	}
	return "user:" + email
}
