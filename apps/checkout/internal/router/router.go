package router

import (
	"net/http"

	"github.com/useportcall/portcall/apps/checkout/internal/modules/address"
	"github.com/useportcall/portcall/apps/checkout/internal/modules/checkout_session"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	r := routerx.New(db, crypto, q)

	// --- Serve static files ---
	r.Use(routerx.StaticFileMiddleware("./frontend/dist"))

	// --- Fallback to index.html for SPA routing ---
	r.NoRoute(func(c *routerx.Context) {
		c.Status(http.StatusOK)
		c.File("./frontend/dist/index.html")
	})

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"message": "pong"}) })

	r.GET("/api/checkout-sessions/:id", checkout_session.GetCheckoutSession)
	r.POST("/api/checkout-sessions/:id/address", address.UpdateCheckoutSessionAddress)

	return r
}
