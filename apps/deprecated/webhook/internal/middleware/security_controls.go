package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/useportcall/portcall/libs/go/routerx"
)

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type limiterStore struct {
	mu      sync.Mutex
	entries map[string]*limiterEntry
	limit   rate.Limit
	burst   int
}

func newLimiterStore(limit rate.Limit, burst int) *limiterStore {
	return &limiterStore{
		entries: make(map[string]*limiterEntry),
		limit:   limit,
		burst:   burst,
	}
}

func (s *limiterStore) get(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.entries[key]; ok {
		entry.lastSeen = time.Now()
		return entry.limiter
	}

	lim := rate.NewLimiter(s.limit, s.burst)
	s.entries[key] = &limiterEntry{limiter: lim, lastSeen: time.Now()}

	if len(s.entries) > 10000 {
		cutoff := time.Now().Add(-30 * time.Minute)
		for k, v := range s.entries {
			if v.lastSeen.Before(cutoff) {
				delete(s.entries, k)
			}
		}
	}

	return lim
}

// RateLimit applies basic IP-based rate limiting.
// Tune with WEBHOOK_RPS/WEBHOOK_BURST and WEBHOOK_STRIPE_RPS/WEBHOOK_STRIPE_BURST.
func RateLimit() routerx.HandlerFunc {
	defaultRPS, defaultBurst := getRateConfig("WEBHOOK_RPS", 10), getBurstConfig("WEBHOOK_BURST", 20)
	stripeRPS, stripeBurst := getRateConfig("WEBHOOK_STRIPE_RPS", 100), getBurstConfig("WEBHOOK_STRIPE_BURST", 200)

	defaultStore := newLimiterStore(defaultRPS, defaultBurst)
	stripeStore := newLimiterStore(stripeRPS, stripeBurst)

	return func(c *routerx.Context) {
		path := c.Request.URL.Path
		clientID := c.ClientIP()
		store := defaultStore
		if strings.HasPrefix(path, "/stripe/") {
			store = stripeStore
		}

		lim := store.get(clientID)
		if !lim.Allow() {
			c.Header("Retry-After", "1")
			c.JSON(http.StatusTooManyRequests, map[string]any{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StripeIPAllowlist limits /stripe/* requests to CIDRs in WEBHOOK_ALLOWED_IPS (optional).
func StripeIPAllowlist() routerx.HandlerFunc {
	cidrs := parseCIDRList(os.Getenv("WEBHOOK_ALLOWED_IPS"))
	if len(cidrs) == 0 {
		return func(c *routerx.Context) { c.Next() }
	}

	return func(c *routerx.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/stripe/") {
			c.Next()
			return
		}

		ip := net.ParseIP(c.ClientIP())
		if ip == nil {
			c.JSON(http.StatusForbidden, map[string]any{"error": "forbidden"})
			c.Abort()
			return
		}

		for _, cidr := range cidrs {
			if cidr.Contains(ip) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, map[string]any{"error": "forbidden"})
		c.Abort()
	}
}

func SecurityHeaders() routerx.HandlerFunc {
	return func(c *routerx.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-site")
		c.Header("Cache-Control", "no-store")

		if isHTTPS(c) {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

func RequestID() routerx.HandlerFunc {
	return func(c *routerx.Context) {
		if c.GetHeader("X-Request-ID") == "" {
			c.Header("X-Request-ID", randomID(16))
		}
		c.Next()
	}
}

func isHTTPS(c *routerx.Context) bool {
	if c.Request.TLS != nil {
		return true
	}

	xfp := strings.ToLower(c.GetHeader("X-Forwarded-Proto"))
	return xfp == "https"
}

func parseCIDRList(raw string) []*net.IPNet {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	var out []*net.IPNet
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		if !strings.Contains(p, "/") {
			p = p + "/32"
		}
		_, cidr, err := net.ParseCIDR(p)
		if err == nil {
			out = append(out, cidr)
		}
	}

	return out
}

func getRateConfig(env string, fallback int) rate.Limit {
	value := strings.TrimSpace(os.Getenv(env))
	if value == "" {
		return rate.Limit(fallback)
	}
	parsed, err := parseInt(value)
	if err != nil || parsed <= 0 {
		return rate.Limit(fallback)
	}
	return rate.Limit(parsed)
}

func getBurstConfig(env string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(env))
	if value == "" {
		return fallback
	}
	parsed, err := parseInt(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseInt(raw string) (int, error) {
	var out int
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, errors.New("non-numeric")
		}
		out = out*10 + int(ch-'0')
	}
	return out, nil
}

func randomID(bytes int) string {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	return hex.EncodeToString(buf)
}
