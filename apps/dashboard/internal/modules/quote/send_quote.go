package quote

import (
	"fmt"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func SendQuote(c *routerx.Context) {
	id := c.Param("id")
	response := new(apix.Quote)

	quoteAppURL := os.Getenv("QUOTE_APP_URL")
	if quoteAppURL == "" {
		c.ServerError("QUOTE_APP_URL is not set", nil)
		return
	}

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}
	if quote.Status == "voided" || quote.Status == "accepted" || quote.Status == "rejected" {
		c.BadRequest("Cannot send quote in current status")
		return
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", c.AppID()); err != nil {
		c.ServerError("Failed to get company", err)
		return
	}

	if quote.UserID == nil {
		c.BadRequest("Quote has no associated user")
		return
	}

	var user models.User
	if err := c.DB().FindForID(*quote.UserID, &user); err != nil {
		c.ServerError("Failed to get user", err)
		return
	}

	recipientEmail := quote.RecipientEmail
	if recipientEmail == "" {
		recipientEmail = user.Email
	}
	if recipientEmail == "" {
		c.BadRequest("Quote has no recipient email")
		return
	}
	quote.RecipientEmail = recipientEmail

	now := time.Now().UTC()
	if quote.Status == "draft" {
		quote.Status = "sent"
		quote.IssuedAt = &now
	}
	if quote.ExpiresAt == nil || quote.ExpiresAt.UTC().Before(now) {
		expiresAt := quoteExpiry(&quote, now)
		quote.ExpiresAt = &expiresAt
	}

	url, err := buildQuoteAccessURL(quoteAppURL, &quote, c.Crypto(), now)
	if err != nil {
		c.ServerError("Failed to build quote access URL", err)
		return
	}
	quote.URL = url

	if err := c.DB().Save(&quote); err != nil {
		c.ServerError("Failed to update quote", err)
		return
	}

	// Trigger send quote email
	if err := c.Queue().Enqueue("send_quote_email", map[string]any{
		"QuoteURL":        *quote.URL,
		"CompanyName":     company.Name,
		"CustomerName":    user.Name,
		"recipient_email": recipientEmail,
	}, "email_queue"); err != nil {
		c.ServerError("Failed to enqueue send quote email job", err)
		return
	}

	if company.Email != "" {
		if err := c.Queue().Enqueue("send_quote_status_email", map[string]any{
			"recipient_email": company.Email,
			"subject":         fmt.Sprintf("Quote issued to %s", recipientEmail),
			"title":           "Quote issued",
			"message":         fmt.Sprintf("Quote %s has been issued to %s.", quote.PublicID, recipientEmail),
			"action_text":     "Open quote",
			"action_url":      *quote.URL,
			"year":            now.Year(),
		}, "email_queue"); err != nil {
			c.ServerError("Failed to enqueue issuer notification", err)
			return
		}
	}

	c.OK(response.Set(&quote))
}
