package apps

import (
	"strconv"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type AppDetail struct {
	ID                uint      `json:"id"`
	PublicID          string    `json:"public_id"`
	Name              string    `json:"name"`
	IsLive            bool      `json:"is_live"`
	AccountID         uint      `json:"account_id"`
	AccountEmail      string    `json:"account_email"`
	UserCount         int64     `json:"user_count"`
	SubscriptionCount int64     `json:"subscription_count"`
	PlanCount         int64     `json:"plan_count"`
	QuoteCount        int64     `json:"quote_count"`
	InvoiceCount      int64     `json:"invoice_count"`
	FeatureCount      int64     `json:"feature_count"`
	ConnectionCount   int64     `json:"connection_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func GetApp(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	// Get account email
	var account models.Account
	accountEmail := ""
	if err := c.DB().FindForID(app.AccountID, &account); err == nil {
		accountEmail = account.Email
	}

	var userCount, subCount, planCount, quoteCount, invoiceCount, featureCount, connectionCount int64
	c.DB().Count(&userCount, models.User{}, "app_id = ?", app.ID)
	c.DB().Count(&subCount, models.Subscription{}, "app_id = ?", app.ID)
	c.DB().Count(&planCount, models.Plan{}, "app_id = ?", app.ID)
	c.DB().Count(&quoteCount, models.Quote{}, "app_id = ?", app.ID)
	c.DB().Count(&invoiceCount, models.Invoice{}, "app_id = ?", app.ID)
	c.DB().Count(&featureCount, models.Feature{}, "app_id = ?", app.ID)
	c.DB().Count(&connectionCount, models.Connection{}, "app_id = ?", app.ID)

	result := AppDetail{
		ID:                app.ID,
		PublicID:          app.PublicID,
		Name:              app.Name,
		IsLive:            app.IsLive,
		AccountID:         app.AccountID,
		AccountEmail:      accountEmail,
		UserCount:         userCount,
		SubscriptionCount: subCount,
		PlanCount:         planCount,
		QuoteCount:        quoteCount,
		InvoiceCount:      invoiceCount,
		FeatureCount:      featureCount,
		ConnectionCount:   connectionCount,
		CreatedAt:         app.CreatedAt,
		UpdatedAt:         app.UpdatedAt,
	}

	c.OK(result)
}
