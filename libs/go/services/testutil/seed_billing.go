//go:build integration

package testutil

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// SeedSubscription creates a Subscription tied to the given App/User/Plan.
func SeedSubscription(
	t *testing.T, db dbx.IORM,
	appID, userID uint, billingAddrID *uint, planID *uint,
) models.Subscription {
	t.Helper()
	now := time.Now()
	sub := models.Subscription{
		PublicID:         dbx.GenPublicID("sub"),
		AppID:            appID,
		UserID:           userID,
		Status:           "active",
		Currency:         "USD",
		InvoiceDueByDays: 10,
		BillingAddressID: billingAddrID,
		PlanID:           planID,
		LastResetAt:      now,
		NextResetAt:      now.AddDate(0, 1, 0),
	}
	if err := db.Create(&sub); err != nil {
		t.Fatalf("seed subscription: %v", err)
	}
	return sub
}

// SeedSubscriptionItem creates a SubscriptionItem for the given Subscription.
func SeedSubscriptionItem(
	t *testing.T, db dbx.IORM,
	appID, subID uint, pricing string, unitAmount int64, qty int32,
) models.SubscriptionItem {
	t.Helper()
	si := models.SubscriptionItem{
		PublicID:       dbx.GenPublicID("si"),
		AppID:          appID,
		SubscriptionID: subID,
		PricingModel:   pricing,
		UnitAmount:     unitAmount,
		Quantity:       qty,
		Title:          "Test Item",
		Description:    "Integration test item",
	}
	if err := db.Create(&si); err != nil {
		t.Fatalf("seed subscription item: %v", err)
	}
	return si
}
