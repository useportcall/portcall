package app

import (
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/apps/quote/internal/handlers"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/libs/go/storex"
)

func NewRouter(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue, store storex.IStore, templateGlob string) routerx.IRouter {
	r := routerx.New(db, crypto, q)
	r.SetStore(store)

	r.Use(func(c *routerx.Context) {
		c.Header("Content-Security-Policy", "frame-ancestors "+frameAncestorsDirective())
		c.Next()
	})

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"status": "ok"}) })

	rateLimiter := ratelimitx.NewGinAdapter()
	r.Use(func(c *routerx.Context) {
		if c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}
		rateLimiter.Middleware(ratelimitx.ByClientIP, 60, time.Minute)(c.Context)
	})

	if strings.TrimSpace(templateGlob) == "" {
		templateGlob = "templates/*.html"
	}
	r.LoadHTMLGlob(templateGlob)
	r.GET("/quotes/:id", handlers.GetQuote)
	r.POST("/quotes/:id", handlers.SubmitQuote)
	r.POST("/quotes/:id/decline", handlers.DeclineQuote)
	r.GET("/quotes/:id/success", handlers.QuoteSuccess)

	return r
}

func frameAncestorsDirective() string {
	origins := []string{"'self'"}
	for _, key := range []string{"DASHBOARD_APP_URL", "DASHBOARD_URL"} {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		origins = append(origins, value)
	}
	if len(origins) == 1 && os.Getenv("APP_ENV") == "development" {
		origins = append(origins, "http://localhost:8080", "http://localhost:8082")
	}
	return strings.Join(origins, " ")
}
