package invoice

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListInvoices(c *routerx.Context) {
	subscriptionID := c.Query("subscription_id")
	userID := c.Query("user_id")

	var conds []any
	if subscriptionID != "" {
		var subscription models.Subscription
		if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
			c.NotFound("Subscription not found")
			return
		}
		conds = []any{"status = 'published' AND app_id = ? AND subscription_id = ?", c.AppID(), subscription.ID}
	} else if userID != "" {
		var user models.User
		if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
			c.NotFound("User not found")
			return
		}

		conds = []any{"app_id = ? AND user_id = ?", c.AppID(), user.ID}
	} else {
		conds = []any{"app_id = ?", c.AppID()}
	}

	var invoices []models.Invoice
	if err := c.DB().ListWithOrder(&invoices, "created_at DESC", conds...); err != nil {
		c.ServerError("Failed to list invoices", err)
		return
	}

	response := make([]apix.Invoice, len(invoices))
	for i, inv := range invoices {
		response[i].Set(&inv)
	}

	c.OK(response)
}
