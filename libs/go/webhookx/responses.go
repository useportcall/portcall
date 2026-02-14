package webhookx

import (
	"net/http"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func respondStripeOK(c *routerx.Context) {
	c.JSON(http.StatusOK, map[string]any{"received": true})
}

func respondStripeError(c *routerx.Context) {
	c.JSON(http.StatusInternalServerError, map[string]any{"error": "internal error"})
}

func respondBraintreeOK(c *routerx.Context) {
	c.JSON(http.StatusOK, map[string]any{"received": true})
}

func respondGenericError(c *routerx.Context, status int) {
	c.JSON(status, map[string]any{"error": "request failed"})
}
