package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeclineQuote(c *routerx.Context) {
	var quote models.Quote
	if err := c.DB().FindFirst(&quote, "public_id = ?", c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return
	}
	if !verifyQuoteAccess(c, &quote) {
		return
	}
	if !quoteCanBeAccepted(&quote, time.Now().UTC()) {
		c.String(http.StatusBadRequest, "quote can no longer be declined")
		return
	}

	now := time.Now().UTC()
	quote.Status = "rejected"
	quote.VoidedAt = &now
	if err := c.DB().Save(&quote); err != nil {
		c.ServerError("failed to save rejected quote", err)
		return
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", quote.AppID); err == nil && company.Email != "" {
		_ = c.Queue().Enqueue("send_quote_status_email", map[string]any{
			"recipient_email": company.Email,
			"subject":         "Quote declined",
			"title":           "Quote declined",
			"message":         fmt.Sprintf("Quote %s has been declined by the recipient.", quote.PublicID),
			"year":            now.Year(),
		}, "email_queue")
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/quotes/%s/success?state=declined", quote.PublicID))
}
