package handlers

import (
	"github.com/useportcall/portcall/libs/go/routerx"
)

// respondGenericError returns a generic error response without leaking details.
// Use this for non-Stripe webhooks that still need a proper error status code.
// For Stripe webhooks, use respondOK instead.
func respondGenericError(c *routerx.Context, status int) {
	c.JSON(status, map[string]any{"error": "request failed"})
}
