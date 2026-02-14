package ratelimitx

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// KeyExtractor extracts a unique key from the request for rate limiting.
type KeyExtractor func(r *http.Request) string

// Middleware creates a rate limiting middleware for HTTP handlers.
func Middleware(limiter *RateLimiter, keyExtractor KeyExtractor, config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyExtractor(r)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			result, err := limiter.Check(key, config)
			if err != nil {
				log.Printf("Rate limit check failed: %v, allowing request", err)
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt.Unix(), 10))

			if !result.Allowed {
				w.Header().Set("Retry-After", strconv.FormatInt(int64(result.ResetAt.Sub(result.ResetAt)*-1), 10))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(fmt.Sprintf(`{"error":"Rate limit exceeded. Try again at %s"}`, result.ResetAt.Format("2006-01-02T15:04:05Z07:00"))))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ByAPIKey extracts the API key from the x-api-key header.
func ByAPIKey(r *http.Request) string {
	apiKey := r.Header.Get("x-api-key")
	if apiKey == "" {
		return ""
	}
	return "apikey:" + apiKey
}

// ByIP extracts the client IP address.
func ByIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return "ip:" + ip
}

// ByUserEmail extracts the user email from context.
func ByUserEmail(contextKey string) KeyExtractor {
	return func(r *http.Request) string {
		email := r.Context().Value(contextKey)
		if email == nil {
			return ""
		}
		emailStr, ok := email.(string)
		if !ok {
			return ""
		}
		return "user:" + emailStr
	}
}

// ByCustom creates a custom key extractor with a prefix.
func ByCustom(prefix string, extractor func(*http.Request) string) KeyExtractor {
	return func(r *http.Request) string {
		value := extractor(r)
		if value == "" {
			return ""
		}
		return prefix + ":" + value
	}
}

