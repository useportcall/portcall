package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type SubscriptionItem struct {
	ID            string     `db:"id" json:"id"`
	Quantity      int32      `db:"quantity" json:"quantity"`
	Interval      string     `json:"interval"`        // billing interval: inherit, week, month, year
	IntervalCount int        `json:"interval_count"`  // number of intervals for the billing cycle
	LastResetAt   *time.Time `json:"last_reset_at"`   // when this item was last billed/reset
	NextResetAt   *time.Time `json:"next_reset_at"`   // when this item will next be billed/reset
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

func (si *SubscriptionItem) Set(subscriptionItem *models.SubscriptionItem) *SubscriptionItem {
	si.ID = subscriptionItem.PublicID
	si.Quantity = subscriptionItem.Quantity
	si.Interval = subscriptionItem.Interval
	si.IntervalCount = subscriptionItem.IntervalCount
	si.LastResetAt = subscriptionItem.LastResetAt
	si.NextResetAt = subscriptionItem.NextResetAt
	si.CreatedAt = subscriptionItem.CreatedAt
	si.UpdatedAt = subscriptionItem.UpdatedAt
	return si
}
