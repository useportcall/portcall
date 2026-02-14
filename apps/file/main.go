package main

import (
	"log"
	"os"
	"time"

	"github.com/useportcall/portcall/apps/file/handlers"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/libs/go/storex"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db, err := dbx.New()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	store, err := storex.New()
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}

	r := routerx.New(db, nil, nil)
	r.SetStore(store)

	// health endpoint
	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"status": "ok"}) })

	// Rate limiting middleware - 120 requests per minute per IP (higher for file serving)
	rateLimiter := ratelimitx.NewGinAdapter()
	r.Use(func(c *routerx.Context) {
		// Skip rate limiting for health check
		if c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}

		// Apply rate limiting by IP
		middleware := rateLimiter.Middleware(
			ratelimitx.ByClientIP,
			120,
			1*time.Minute,
		)
		middleware(c.Context)
	})

	r.GET("/invoice/:invoice_id", handlers.GetInvoice)
	r.GET("/invoices/:invoice_id/view", handlers.ViewInvoice)
	r.GET("/icon-logos/:filename", handlers.GetIconLogo)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
