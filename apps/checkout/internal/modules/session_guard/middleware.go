package session_guard

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// SessionHandler is a handler that receives a pre-authenticated checkout session.
type SessionHandler func(c *routerx.Context, session *models.CheckoutSession)

// WithCheckoutSession wraps a handler with checkout session authentication.
// It validates the session token, loads the session from the DB, checks
// status/expiry, and passes the authenticated session to the inner handler.
// This eliminates the need for each handler to call AuthorizeCheckoutSession.
func WithCheckoutSession(handler SessionHandler) routerx.HandlerFunc {
	return func(c *routerx.Context) {
		session, ok := AuthorizeCheckoutSession(c)
		if !ok {
			return
		}
		handler(c, session)
	}
}
