package quote

import (
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateQuoteRequest struct {
	UserID                string     `json:"user_id"`
	ExpiresAt             *time.Time `json:"expires_at"`
	CompanyName           *string    `json:"company_name"`
	DirectCheckoutEnabled *bool      `json:"direct_checkout_enabled"`
	Toc                   *string    `json:"toc"`
	RecipientEmail        *string    `json:"recipient_email"`
	RecipientName         *string    `json:"recipient_name"`
	RecipientTitle        *string    `json:"recipient_title"`
	PreparedByEmail       *string    `json:"prepared_by_email"`
}

func UpdateQuote(c *routerx.Context) {
	id := c.Param("id")

	var body UpdateQuoteRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}
	if !IsMutableStatus(quote.Status) {
		c.BadRequest("Quote can no longer be edited")
		return
	}

	if body.UserID != "" {
		var user models.User
		if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
			c.NotFound("User not found")
			return
		}

		quote.UserID = &user.ID
	}

	if body.CompanyName != nil {
		quote.CompanyName = *body.CompanyName
	}

	if body.DirectCheckoutEnabled != nil {
		quote.DirectCheckout = *body.DirectCheckoutEnabled
	}

	if body.Toc != nil {
		quote.Toc = *body.Toc
	}

	if body.RecipientEmail != nil {
		quote.RecipientEmail = *body.RecipientEmail
	}

	if body.RecipientName != nil {
		quote.PublicName = *body.RecipientName
	}

	if body.RecipientTitle != nil {
		quote.PublicTitle = *body.RecipientTitle
	}

	if body.PreparedByEmail != nil {
		quote.PreparedByEmail = *body.PreparedByEmail
	}

	quote.ExpiresAt = body.ExpiresAt

	if err := c.DB().Save(&quote); err != nil {
		c.ServerError("Failed to update quote", err)
		return
	}

	response := new(apix.Quote)

	c.OK(response.Set(&quote))
}
