package handlers

import (
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func verifyQuoteAccess(c *routerx.Context, quote *models.Quote) bool {
	if quote.Status == "draft" {
		c.String(403, "quote is not yet available")
		return false
	}
	token := c.Query("qt")
	if err := cryptox.VerifyQuoteAccessToken(c.Crypto(), token, quote.PublicID, time.Now().UTC()); err != nil {
		c.String(403, "invalid or expired quote link")
		return false
	}
	return true
}

func quoteCanBeAccepted(quote *models.Quote, now time.Time) bool {
	if quote.Status != "sent" {
		return false
	}
	if quote.ExpiresAt == nil {
		return true
	}
	return now.UTC().Before(quote.ExpiresAt.UTC()) || now.UTC().Equal(quote.ExpiresAt.UTC())
}

func quoteStatusBanner(status string) string {
	switch status {
	case "accepted":
		return "This quote has already been accepted."
	case "voided":
		return "This quote has been withdrawn by the issuer."
	case "rejected":
		return "This quote has been declined."
	default:
		return ""
	}
}
