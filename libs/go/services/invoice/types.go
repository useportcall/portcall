package invoice

import "github.com/useportcall/portcall/libs/go/dbx/models"

// ListInput is the input for listing invoices.
type ListInput struct {
	AppID          uint
	SubscriptionID string // optional public ID filter
	UserID         string // optional public ID filter
}

// ListResult is the result of listing invoices.
type ListResult struct {
	Invoices []models.Invoice
}

// CreateInput is the input for creating an invoice with items.
type CreateInput struct {
	SubscriptionID uint
}

// CreateResult is the result of creating an invoice.
type CreateResult struct {
	Invoice *models.Invoice
	Skipped bool
}

// CreateUpgradeInput is the input for creating an upgrade invoice.
type CreateUpgradeInput struct {
	SubscriptionID  uint  `json:"subscription_id"`
	PriceDifference int64 `json:"price_difference"`
	OldPlanID       uint  `json:"old_plan_id"`
	NewPlanID       uint  `json:"new_plan_id"`
}

// CreateUpgradeResult is the result of creating an upgrade invoice.
type CreateUpgradeResult struct {
	Invoice   *models.Invoice
	ShouldPay bool
}
