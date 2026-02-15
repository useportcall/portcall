package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/webhook/internal/handlers"
	"github.com/useportcall/portcall/apps/webhook/internal/middleware"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT environment variable not set")
	}

	db, err := dbx.New()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := qx.New()
	if err != nil {
		log.Fatal(err)
	}
	crypto, err := cryptox.New()
	if err != nil {
		log.Fatal(err)
	}

	r := routerx.New(db, crypto, queue)

	// security + traffic controls
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RequestID())
	r.Use(middleware.RateLimit())
	r.Use(middleware.StripeIPAllowlist())

	// health endpoints
	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"status": "ok"}) })
	r.GET("/healthz", func(c *routerx.Context) { c.OK(map[string]any{"status": "ok"}) })

	r.POST("/stripe/:connection_id", handlers.HandleStripeWebhook)
	r.POST("/postmark/:connection_id", handlers.HandlePostmarkWebhook)

	// default 404 for all other routes
	r.NoRoute(func(c *routerx.Context) { c.NotFound("Not found") })

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
