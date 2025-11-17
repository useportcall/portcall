package main

import (
	"log"
	"os"
	"strings"

	"github.com/useportcall/portcall/apps/api/internal/modules/v1/checkout_session"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/entitlement"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/meter_event"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/subscription"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/user"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db := dbx.New()
	crypto := cryptox.New()
	q := qx.New()

	r := routerx.New(db, crypto, q)

	r.Use(func(c *routerx.Context) {
		apiKey := c.Request.Header.Get("x-api-key")
		if apiKey == "" {
			log.Println("Missing API key")
			c.Unauthorized("Missing API key")
			c.Abort()
			return
		}

		parts := strings.Split(apiKey, "_")
		publicID := strings.Join(parts[0:2], "_")
		sk := parts[2]

		var secret models.Secret
		if err := c.DB().FindFirst(&secret, "public_id = ?", publicID); err != nil {
			log.Println("API Key lookup failed:", err)
			c.Unauthorized("Invalid API key")
			c.Abort()
			return
		}

		if secret.DisabledAt != nil {
			log.Println("API Key is disabled")
			c.Unauthorized("Invalid API key")
			c.Abort()
			return
		}

		ok, err := c.Crypto().CompareHash(secret.KeyHash, sk)
		if err != nil || !ok {
			log.Println("API Key comparison failed:", err)
			c.Unauthorized("Invalid API key")
			c.Abort()
			return
		}

		c.Set("app_id", secret.AppID)

		c.Next()
	})

	r.GET("/v1/users/", user.ListUsers)

	r.GET("/v1/users/:user_id/entitlements/:id", entitlement.GetEntitlement)

	r.GET("/v1/subscriptions", subscription.ListSubscriptions)
	r.POST("/v1/subscriptions/:subscription_id", subscription.UpdateSubscription)
	r.POST("/v1/subscriptions/:subscription_id/cancel", subscription.CancelSubscription)

	r.POST("/v1/meter-events", meter_event.CreateMeterEvent)

	r.POST("/v1/checkout-sessions", checkout_session.CreateCheckoutSession)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
