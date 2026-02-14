package webhookx

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func (w *Router) guard(c *routerx.Context, provider string) bool {
	setSecurityHeaders(c)
	ensureRequestID(c)
	if !w.allowRate(c, provider) {
		c.Header("Retry-After", "1")
		c.JSON(http.StatusTooManyRequests, map[string]any{"error": "rate limit exceeded"})
		return false
	}
	if provider == "stripe" && !w.allowStripeIP(c) {
		c.JSON(http.StatusForbidden, map[string]any{"error": "forbidden"})
		return false
	}
	return true
}

func (w *Router) allowRate(c *routerx.Context, provider string) bool {
	store := w.defaultStore
	switch provider {
	case "stripe":
		store = w.stripeStore
	case "braintree":
		store = w.braintreeStore
	}
	return store.allow(c.ClientIP())
}

func (w *Router) allowStripeIP(c *routerx.Context) bool {
	if len(w.stripeCIDRs) == 0 {
		return true
	}
	ip := net.ParseIP(c.ClientIP())
	if ip == nil {
		return false
	}
	for _, cidr := range w.stripeCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func ensureRequestID(c *routerx.Context) {
	if c.GetHeader("X-Request-ID") != "" {
		return
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err == nil {
		c.Header("X-Request-ID", hex.EncodeToString(buf))
	}
}

func setSecurityHeaders(c *routerx.Context) {
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
}

func isHTTPS(c *routerx.Context) bool {
	if c.Request.TLS != nil {
		return true
	}
	return strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
}
