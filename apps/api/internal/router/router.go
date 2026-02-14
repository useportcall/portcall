package router

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// Init creates the public API router with API-key auth and rate limiting.
func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	r := routerx.New(db, crypto, q)

	r.GET("/ping", func(c *routerx.Context) {
		c.OK(map[string]any{"status": "ok"})
	})
	r.GET("/healthz", func(c *routerx.Context) {
		c.OK(map[string]any{"status": "healthy"})
	})

	r.Use(apiKeyAuth)
	r.Use(rateLimitMiddleware)

	registerRoutes(r)
	return r
}

func apiKeyAuth(c *routerx.Context) {
	if c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/healthz" {
		c.Next()
		return
	}

	apiKey := c.Request.Header.Get("x-api-key")
	if apiKey == "" {
		c.Unauthorized("Missing API key")
		c.Abort()
		return
	}

	parts := strings.Split(apiKey, "_")
	if len(parts) < 3 {
		c.Unauthorized("Invalid API key format")
		c.Abort()
		return
	}
	publicID := strings.Join(parts[0:2], "_")
	sk := parts[2]

	var secret models.Secret
	if err := c.DB().FindFirst(&secret, "public_id = ?", publicID); err != nil {
		c.Unauthorized("Invalid API key")
		c.Abort()
		return
	}
	if secret.DisabledAt != nil {
		c.Unauthorized("Invalid API key")
		c.Abort()
		return
	}

	ok, err := c.Crypto().CompareHash(secret.KeyHash, sk)
	if err != nil || !ok {
		c.Unauthorized("Invalid API key")
		c.Abort()
		return
	}

	var app models.App
	if err := c.DB().FindForID(secret.AppID, &app); err != nil {
		c.Unauthorized("Invalid API key")
		c.Abort()
		return
	}
	c.Set("app_id", app.ID)
	c.Set("public_app_id", app.PublicID)
	c.Set("is_live", app.IsLive)
	c.Set("is_billing_exempt", app.BillingExempt)
	c.Next()
}
