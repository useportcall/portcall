package router

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/apps/checkout/internal/modules/address"
	"github.com/useportcall/portcall/apps/checkout/internal/modules/checkout_session"
	"github.com/useportcall/portcall/apps/checkout/internal/modules/payment_link"
	"github.com/useportcall/portcall/apps/checkout/internal/modules/session_guard"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	r := routerx.New(db, crypto, q)

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"message": "pong"}) })
	r.GET("/healthz", func(c *routerx.Context) { c.OK(map[string]any{"status": "healthy"}) })

	r.Use(func(c *routerx.Context) {
		headers := c.Writer.Header()
		headers.Set("X-Content-Type-Options", "nosniff")
		headers.Set("X-Frame-Options", "DENY")
		headers.Set("Referrer-Policy", "strict-origin")
		headers.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		headers.Set("Cross-Origin-Resource-Policy", "same-origin")
		headers.Set("Content-Security-Policy", buildCheckoutCSP())
		if c.Request.TLS != nil {
			headers.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			headers.Set("Cache-Control", "no-store")
		}
		c.Next()
	})

	// --- Serve static files ---
	staticDir := "./frontend/dist"
	if v := os.Getenv("CHECKOUT_STATIC_DIR"); v != "" {
		staticDir = v
	}
	r.Use(routerx.StaticFileMiddleware(staticDir))

	// Rate limiting middleware - 60 requests per minute per IP for API endpoints
	rateLimiter := ratelimitx.NewGinAdapter()
	r.Use(func(c *routerx.Context) {
		// Only rate limit API endpoints
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}

		// Apply rate limiting by IP
		middleware := rateLimiter.Middleware(
			ratelimitx.ByClientIP,
			60,
			1*time.Minute,
		)
		middleware(c.Context)
	})

	// --- Fallback: 404 for unknown API routes, SPA for everything else ---
	r.NoRoute(func(c *routerx.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		c.Status(http.StatusOK)
		c.File(staticDir + "/index.html")
	})

	// Checkout-session routes use the session auth middleware.
	auth := session_guard.WithCheckoutSession
	r.GET("/api/checkout-sessions/:id", auth(checkout_session.GetCheckoutSession))
	r.POST("/api/checkout-sessions/:id/address", auth(address.UpdateCheckoutSessionAddress))
	r.POST("/api/checkout-sessions/:id/complete", auth(checkout_session.CompleteCheckoutSession))

	// Payment link route uses its own token-based auth.
	r.POST("/api/payment-links/:id/redeem", payment_link.RedeemPaymentLink)

	return r
}
