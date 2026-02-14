package subscription

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// FindResetsResult is the result of finding subscriptions to reset.
type FindResetsResult struct {
	SubscriptionIDs []uint
}

// StartResetInput is the input for starting a subscription reset.
type StartResetInput struct {
	SubscriptionID uint `json:"subscription_id"`
}

// StartResetResult is the result of starting a subscription reset.
type StartResetResult struct {
	Subscription   *models.Subscription
	User           *models.User
	Status         string // "active", "canceled", or "skipped"
	CheckoutURL    string
	RollbackPlanID *uint
	FinalResetAt   time.Time
}

// EndResetInput is the input for ending a subscription reset.
type EndResetInput struct {
	SubscriptionID uint `json:"subscription_id"`
}

// EndResetResult is the result of ending a subscription reset.
type EndResetResult struct {
	Subscription  *models.Subscription
	LastResetAt   time.Time
	NextResetAt   time.Time
	AppliedPlanID *uint
}
