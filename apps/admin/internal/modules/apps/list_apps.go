package apps

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type AppListItem struct {
	ID                uint      `json:"id"`
	PublicID          string    `json:"public_id"`
	Name              string    `json:"name"`
	IsLive            bool      `json:"is_live"`
	UserCount         int64     `json:"user_count"`
	SubscriptionCount int64     `json:"subscription_count"`
	PlanCount         int64     `json:"plan_count"`
	QuoteCount        int64     `json:"quote_count"`
	InvoiceCount      int64     `json:"invoice_count"`
	CreatedAt         time.Time `json:"created_at"`
}

func ListApps(c *routerx.Context) {
	var apps []models.App
	err := c.DB().ListWithOrder(&apps, "created_at DESC")
	if err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	result := make([]AppListItem, len(apps))
	for i, app := range apps {
		var userCount, subCount, planCount, quoteCount, invoiceCount int64
		c.DB().Count(&userCount, models.User{}, "app_id = ?", app.ID)
		c.DB().Count(&subCount, models.Subscription{}, "app_id = ?", app.ID)
		c.DB().Count(&planCount, models.Plan{}, "app_id = ?", app.ID)
		c.DB().Count(&quoteCount, models.Quote{}, "app_id = ?", app.ID)
		c.DB().Count(&invoiceCount, models.Invoice{}, "app_id = ?", app.ID)

		result[i] = AppListItem{
			ID:                app.ID,
			PublicID:          app.PublicID,
			Name:              app.Name,
			IsLive:            app.IsLive,
			UserCount:         userCount,
			SubscriptionCount: subCount,
			PlanCount:         planCount,
			QuoteCount:        quoteCount,
			InvoiceCount:      invoiceCount,
			CreatedAt:         app.CreatedAt,
		}
	}

	c.OK(result)
}
