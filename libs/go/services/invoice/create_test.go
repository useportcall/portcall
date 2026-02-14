package invoice_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
)

func TestCreate_Success(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test.com")
	defer os.Unsetenv("INVOICE_APP_URL")

	billingAddr := uint(50)
	sub := &models.Subscription{
		AppID: 1, UserID: 2, Currency: "USD",
		InvoiceDueByDays: 10, BillingAddressID: &billingAddr,
	}
	sub.ID = 100

	db := &mockDB{
		subscription: sub,
		company:      &models.Company{BillingAddressID: 60},
		invoiceCount: 5,
		subItemIDs:   []uint{10, 20},
		subscriptionItems: map[uint]*models.SubscriptionItem{
			10: {AppID: 1, Quantity: 1, PricingModel: "fixed", UnitAmount: 1000, Title: "Item A"},
			20: {AppID: 1, Quantity: 2, PricingModel: "fixed", UnitAmount: 500, Title: "Item B"},
		},
		notFoundInvoice: true,
	}

	svc := inv.NewService(db)
	result, err := svc.Create(&inv.CreateInput{SubscriptionID: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Skipped {
		t.Fatal("expected created, not skipped")
	}
	if result.Invoice.InvoiceNumber != "INV-0000006" {
		t.Fatalf("expected INV-0000006, got %s", result.Invoice.InvoiceNumber)
	}
	if result.Invoice.Total != 2000 {
		t.Fatalf("expected total 2000, got %d", result.Invoice.Total)
	}
	if result.Invoice.Status != "issued" {
		t.Fatalf("expected status issued, got %s", result.Invoice.Status)
	}
}

func TestCreate_Idempotent_Skips(t *testing.T) {
	sub := &models.Subscription{AppID: 1}
	sub.ID = 100

	db := &mockDB{
		subscription:    sub,
		notFoundInvoice: false, // invoice already exists
	}

	svc := inv.NewService(db)
	result, err := svc.Create(&inv.CreateInput{SubscriptionID: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Skipped {
		t.Fatal("expected skipped for idempotent call")
	}
}
