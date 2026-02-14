package session_guard

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

const (
	checkoutSessionTokenHeader = "X-Checkout-Session-Token"
)

func AuthorizeCheckoutSession(c *routerx.Context) (*models.CheckoutSession, bool) {
	sessionID := strings.TrimSpace(c.Param("id"))
	if !cryptox.IsValidCheckoutSessionID(sessionID) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "invalid or expired checkout session"})
		return nil, false
	}

	token := strings.TrimSpace(c.GetHeader(checkoutSessionTokenHeader))
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "missing checkout session token"})
		return nil, false
	}

	if err := cryptox.VerifyCheckoutSessionToken(c.Crypto(), token, sessionID, time.Now()); err != nil {
		if errors.Is(err, cryptox.ErrExpiredCheckoutSessionToken) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "checkout session has expired"})
			return nil, false
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "invalid checkout session token"})
		return nil, false
	}

	var session models.CheckoutSession
	if err := c.DB().FindFirst(&session, "public_id = ?", sessionID); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "invalid or expired checkout session"})
		return nil, false
	}

	now := time.Now()
	if session.Status != "active" || now.After(session.ExpiresAt) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": "invalid or expired checkout session"})
		return nil, false
	}

	return &session, true
}
