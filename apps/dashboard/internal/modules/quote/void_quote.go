package quote

import (
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func VoidQuote(c *routerx.Context) {
	id := c.Param("id")

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}
	if quote.Status == "voided" || quote.Status == "accepted" || quote.Status == "rejected" {
		c.BadRequest("Quote cannot be voided in current status")
		return
	}

	now := time.Now()
	quote.VoidedAt = &now
	quote.Status = "voided"

	if err := c.DB().Save(&quote); err != nil {
		c.ServerError("Failed to update quote", err)
		return
	}

	if quote.RecipientEmail != "" {
		actionURL := ""
		if quote.URL != nil {
			actionURL = *quote.URL
		}
		_ = c.Queue().Enqueue("send_quote_status_email", map[string]any{
			"recipient_email": quote.RecipientEmail,
			"subject":         "Quote withdrawn",
			"title":           "Quote withdrawn",
			"message":         "This quote has been withdrawn by the issuer.",
			"action_text":     "View quote",
			"action_url":      actionURL,
			"year":            now.Year(),
		}, "email_queue")
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", c.AppID()); err == nil && company.Email != "" {
		_ = c.Queue().Enqueue("send_quote_status_email", map[string]any{
			"recipient_email": company.Email,
			"subject":         "Quote voided",
			"title":           "Quote voided",
			"message":         fmt.Sprintf("Quote %s has been voided.", quote.PublicID),
			"action_text":     "Open dashboard",
			"action_url":      "",
			"year":            now.Year(),
		}, "email_queue")
	}

	response := new(apix.Quote)

	c.OK(response.Set(&quote))
}
