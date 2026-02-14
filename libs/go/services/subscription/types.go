package subscription

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// FindInput is the input for finding a subscription.
type FindInput struct {
	AppID  uint `json:"app_id"`
	UserID uint `json:"user_id"`
	PlanID uint `json:"plan_id"`
}

// FindResult is the result of finding a subscription.
type FindResult struct {
	Action         string // "create" or "update"
	SubscriptionID uint
	AppID          uint
	UserID         uint
	PlanID         uint
}

// CreateInput is the input for creating a subscription.
type CreateInput struct {
	AppID  uint `json:"app_id"`
	PlanID uint `json:"plan_id"`
	UserID uint `json:"user_id"`
}

// CreateResult is the result of creating a subscription.
type CreateResult struct {
	Subscription *models.Subscription
	UserID       uint
	PlanID       uint
	ItemCount    int
}

// UpdateInput is the input for updating a subscription.
type UpdateInput struct {
	SubscriptionID uint `json:"subscription_id"`
	PlanID         uint `json:"plan_id"`
	AppID          uint `json:"app_id"`
}

// UpdateResult is the result of updating a subscription.
type UpdateResult struct {
	Subscription *models.Subscription
	OldPlanID    uint
	NewPlanID    uint
}

// PlanSwitchInput is the input for processing a plan switch.
type PlanSwitchInput struct {
	OldPlanID      uint `json:"old_plan_id"`
	NewPlanID      uint `json:"new_plan_id"`
	SubscriptionID uint `json:"subscription_id"`
}

// PlanSwitchResult is the result of processing a plan switch.
type PlanSwitchResult struct {
	IsUpgrade       bool
	SubscriptionID  uint
	UserID          uint
	OldPlanID       uint
	NewPlanID       uint
	PriceDifference int64
}

// CreateItemsInput is the input for bulk creating items.
type CreateItemsInput struct {
	SubscriptionID uint `json:"subscription_id"`
	PlanID         uint `json:"plan_id"`
}

// CreateItemsResult is the result of bulk creating items.
type CreateItemsResult struct {
	SubscriptionID uint
	ItemCount      int
}

// Prorate calculates the prorated amount based on the fraction
// of the billing period remaining.
func Prorate(amount int64, start, end time.Time, now time.Time) int64 {
	total := end.Sub(start)
	remaining := end.Sub(now)
	if total <= 0 || remaining <= 0 {
		return amount
	}
	if remaining > total {
		return amount
	}
	return amount * int64(remaining) / int64(total)
}
