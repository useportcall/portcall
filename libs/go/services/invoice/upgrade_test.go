package invoice_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
)

func TestCreateUpgrade_Success(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	billingAddr := uint(99)
	sub := &models.Subscription{
		AppID: 1, UserID: 2, Currency: "usd",
		InvoiceDueByDays: 30, BillingAddressID: &billingAddr,
	}
	sub.ID = 10

	db := &mockDB{
		subscription: sub,
		company:      &models.Company{BillingAddressID: 50},
		invoiceCount: 5,
		plans: map[uint]*models.Plan{
			1: {Name: "Basic"},
			2: {Name: "Pro"},
		},
	}

	svc := inv.NewService(db)
	result, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{
		SubscriptionID: 10, PriceDifference: 1500, OldPlanID: 1, NewPlanID: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.ShouldPay {
		t.Fatal("expected ShouldPay = true")
	}
	if result.Invoice.Total != 1500 {
		t.Fatalf("expected total 1500, got %d", result.Invoice.Total)
	}
	if result.Invoice.Status != "issued" {
		t.Fatalf("expected status issued, got %s", result.Invoice.Status)
	}
}

func TestCreateUpgrade_NoURL(t *testing.T) {
	os.Unsetenv("INVOICE_APP_URL")

	db := &mockDB{
		subscription: &models.Subscription{},
	}

	svc := inv.NewService(db)
	_, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{SubscriptionID: 1})
	if err == nil {
		t.Fatal("expected error for missing INVOICE_APP_URL")
	}
}
