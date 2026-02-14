package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Subscription struct {
	ID                    string             `json:"id"`
	AppID                 uint               `json:"app_id"`
	UserID                string             `json:"user_id"`
	PlanID                string             `json:"plan_id"`
	ScheduledPlanID       *string            `json:"scheduled_plan_id,omitempty"`
	Status                string             `json:"status"`
	IsFree                bool               `json:"is_free"`
	LastResetAt           time.Time          `json:"last_reset_at"`
	NextResetAt           time.Time          `json:"next_reset_at"`
	StripePaymentMethodID *string            `json:"stripe_payment_method_id"`
	TrialDurationDays     int                `json:"trial_duration_days"`
	InvoiceCount          int                `json:"invoice_count"` // Number of invoices generated for this subscription
	AutoCollection        bool               `json:"auto_collection"`
	Currency              string             `json:"currency"`
	BillingInterval       string             `json:"billing_interval"` // Consider using a struct for strong typing
	BillingIntervalCount  int                `json:"billing_interval_count"`
	Items                 []SubscriptionItem `json:"items"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
	User                  any                `json:"user"`
	Plan                  any                `json:"plan"`
	ScheduledPlan         any                `json:"scheduled_plan,omitempty"`
}

func (s *Subscription) Set(subscription *models.Subscription) *Subscription {
	s.ID = subscription.PublicID
	s.AppID = subscription.AppID
	s.Status = subscription.Status
	s.IsFree = subscription.IsFree
	s.LastResetAt = subscription.LastResetAt
	s.NextResetAt = subscription.NextResetAt
	s.BillingInterval = subscription.BillingInterval
	s.BillingIntervalCount = subscription.BillingIntervalCount
	s.CreatedAt = subscription.CreatedAt
	s.UpdatedAt = subscription.UpdatedAt
	return s
}
