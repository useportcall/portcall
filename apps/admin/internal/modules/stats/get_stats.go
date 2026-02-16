package stats

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type StatsResponse struct {
	TotalApps           int64 `json:"total_apps"`
	TotalUsers          int64 `json:"total_users"`
	TotalSubscriptions  int64 `json:"total_subscriptions"`
	TotalPlans          int64 `json:"total_plans"`
	TotalQuotes         int64 `json:"total_quotes"`
	TotalInvoices       int64 `json:"total_invoices"`
	TotalAccounts       int64 `json:"total_accounts"`
	ActiveSubscriptions int64 `json:"active_subscriptions"`
	DraftQuotes         int64 `json:"draft_quotes"`
	PaidInvoices        int64 `json:"paid_invoices"`
	PendingInvoices     int64 `json:"pending_invoices"`
}

func GetStats(c *routerx.Context) {
	response := StatsResponse{}

	// Total counts - need empty condition string for all records
	c.DB().Count(&response.TotalApps, models.App{}, "1=1")
	c.DB().Count(&response.TotalUsers, models.User{}, "1=1")
	c.DB().Count(&response.TotalSubscriptions, models.Subscription{}, "1=1")
	c.DB().Count(&response.TotalPlans, models.Plan{}, "1=1")
	c.DB().Count(&response.TotalQuotes, models.Quote{}, "1=1")
	c.DB().Count(&response.TotalInvoices, models.Invoice{}, "1=1")
	c.DB().Count(&response.TotalAccounts, models.Account{}, "1=1")

	// Status-specific counts
	c.DB().Count(&response.ActiveSubscriptions, models.Subscription{}, "status = ?", "active")
	c.DB().Count(&response.DraftQuotes, models.Quote{}, "status = ?", "draft")
	c.DB().Count(&response.PaidInvoices, models.Invoice{}, "status = ?", "paid")
	c.DB().Count(&response.PendingInvoices, models.Invoice{}, "status = ?", "draft")

	c.OK(response)
}
