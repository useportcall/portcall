package subscription

import (
	"time"

	"github.com/useportcall/portcall/apps/dashboard/internal/modules/subscription_item"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Subscription struct {
	ID                    string                               `json:"id"`
	AppID                 uint                                 `json:"app_id"`
	UserID                string                               `json:"user_id"`
	PlanID                string                               `json:"plan_id"`
	Status                string                               `json:"status"`
	LastResetAt           string                               `json:"last_reset_at"`
	NextResetAt           string                               `json:"next_reset_at"`
	StripePaymentMethodID *string                              `json:"stripe_payment_method_id"`
	TrialDurationDays     int                                  `json:"trial_duration_days"`
	InvoiceCount          int                                  `json:"invoice_count"` // Number of invoices generated for this subscription
	AutoCollection        bool                                 `json:"auto_collection"`
	Currency              string                               `json:"currency"`
	BillingInterval       string                               `json:"billing_interval"` // Consider using a struct for strong typing
	BillingIntervalCount  int                                  `json:"billing_interval_count"`
	Items                 []subscription_item.SubscriptionItem `json:"items"`
	CreatedAt             time.Time                            `json:"created_at"`
	UpdatedAt             time.Time                            `json:"updated_at"`
	User                  any                                  `json:"user"`
	Plan                  any                                  `json:"plan"`
}

func (s *Subscription) Set(subscription *models.Subscription) *Subscription {
	s.ID = subscription.PublicID
	s.AppID = subscription.AppID
	s.Status = subscription.Status
	s.LastResetAt = subscription.LastResetAt.Format("2006-01-02T15:04:05Z07:00")
	s.NextResetAt = subscription.NextResetAt.Format("2006-01-02T15:04:05Z07:00")
	s.BillingInterval = subscription.BillingInterval
	s.BillingIntervalCount = subscription.BillingIntervalCount
	s.CreatedAt = subscription.CreatedAt
	s.UpdatedAt = subscription.UpdatedAt
	return s
}
