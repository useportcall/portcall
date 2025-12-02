package quote

import (
	"fmt"
	"os"

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

	if quote.Status == "voided" {
		c.BadRequest("Cannot send a voided quote")
		return
	}

	if quote.Status == "draft" {
		// Update the quote status to "sent"
		quote.Status = "sent"
		if err := c.DB().Save(&quote); err != nil {
			c.ServerError("Failed to update quote status", err)
			return
		}
	}

	// Trigger send quote email
	if err := c.Queue().Enqueue("send_quote_email", map[string]any{
		"QuoteURL":     fmt.Sprintf("%s/quotes/%s", quoteAppURL, quote.PublicID),
		"CompanyName":  company.Name,
		"CustomerName": user.Name,
	}, "email_queue"); err != nil {
		c.ServerError("Failed to enqueue send quote email job", err)
		return
	}

	c.OK(response.Set(&quote))
}
