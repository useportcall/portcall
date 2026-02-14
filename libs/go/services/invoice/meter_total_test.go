package invoice_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
)

func TestCreate_UsesBillingMeterForUnitPricing(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test.com")
	defer os.Unsetenv("INVOICE_APP_URL")

	addr := uint(50)
	planItemID := uint(10)
	sub := &models.Subscription{AppID: 1, UserID: 2, Currency: "USD", BillingAddressID: &addr}
	sub.ID = 100
	db := &mockDB{
		subscription:    sub,
		company:         &models.Company{BillingAddressID: 60},
		subItemIDs:      []uint{1},
		notFoundInvoice: true,
		subscriptionItems: map[uint]*models.SubscriptionItem{
			1: {AppID: 1, PlanItemID: &planItemID, PricingModel: "unit", UnitAmount: 999, Quantity: 1, Usage: 1},
		},
		billingMeters: []models.BillingMeter{
			{SubscriptionID: 100, PlanItemID: 10, PricingModel: "unit", UnitAmount: 99, Usage: 120, FreeQuota: 20},
		},
	}
	r, err := inv.NewService(db).Create(&inv.CreateInput{SubscriptionID: 100})
	if err != nil {
		t.Fatal(err)
	}
	if r.Invoice.Total != 9900 {
		t.Fatalf("expected 9900, got %d", r.Invoice.Total)
	}
}

func TestCreate_UsesGraduatedTieredMeterPricing(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test.com")
	defer os.Unsetenv("INVOICE_APP_URL")

	addr := uint(50)
	planItemID := uint(20)
	sub := &models.Subscription{AppID: 1, UserID: 2, Currency: "USD", BillingAddressID: &addr}
	sub.ID = 200
	tiers := []models.Tier{{Start: 0, End: 100, Amount: 10}, {Start: 100, End: 300, Amount: 5}}
	db := &mockDB{
		subscription:    sub,
		company:         &models.Company{BillingAddressID: 60},
		subItemIDs:      []uint{2},
		notFoundInvoice: true,
		subscriptionItems: map[uint]*models.SubscriptionItem{
			2: {AppID: 1, PlanItemID: &planItemID, PricingModel: "tiered", Quantity: 1, Tiers: &tiers},
		},
		billingMeters: []models.BillingMeter{
			{SubscriptionID: 200, PlanItemID: 20, PricingModel: "tiered", Usage: 250, FreeQuota: 50, Tiers: &tiers},
		},
	}
	r, err := inv.NewService(db).Create(&inv.CreateInput{SubscriptionID: 200})
	if err != nil {
		t.Fatal(err)
	}
	if r.Invoice.Total != 1500 {
		t.Fatalf("expected 1500, got %d", r.Invoice.Total)
	}
}
