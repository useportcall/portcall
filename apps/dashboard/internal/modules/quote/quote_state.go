package quote

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func IsMutableStatus(status string) bool {
	return status == "draft"
}

func HasLockedQuoteForPlan(c *routerx.Context, planID uint) (bool, error) {
	var result models.Quote
	if err := c.DB().FindFirstOrNil(&result,
		"app_id = ? AND plan_id = ? AND status <> ?", c.AppID(), planID, "draft"); err != nil {
		return false, err
	}
	return result.ID != 0, nil
}

func HasLockedQuoteForAddress(c *routerx.Context, addressID uint) (bool, error) {
	var result models.Quote
	if err := c.DB().FindFirstOrNil(&result,
		"app_id = ? AND recipient_address_id = ? AND status <> ?",
		c.AppID(), addressID, "draft"); err != nil {
		return false, err
	}
	return result.ID != 0, nil
}

func quoteExpiry(quote *models.Quote, now time.Time) time.Time {
	if quote.Status == "accepted" {
		return now.UTC().Add(365 * 24 * time.Hour)
	}
	if quote.ExpiresAt != nil && !quote.ExpiresAt.IsZero() {
		if quote.ExpiresAt.UTC().After(now.UTC()) {
			return quote.ExpiresAt.UTC()
		}
		if quote.Status != "sent" {
			return now.UTC().Add(90 * 24 * time.Hour)
		}
		return quote.ExpiresAt.UTC()
	}
	days := quote.DaysValid
	if days <= 0 {
		days = 30
	}
	return now.UTC().Add(time.Duration(days) * 24 * time.Hour)
}
