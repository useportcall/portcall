//go:build integration

package invoice_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	inv "github.com/useportcall/portcall/libs/go/services/invoice"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_Integration_NoSubItems_Errors(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	svc := inv.NewService(testDB)

	_, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	if err == nil {
		t.Fatal("expected error for subscription without items")
	}
}

func TestCreate_Integration_MissingSub_Errors(t *testing.T) {
	svc := inv.NewService(testDB)

	_, err := svc.Create(&inv.CreateInput{SubscriptionID: 999999})
	if err == nil {
		t.Fatal("expected error for nonexistent subscription")
	}
}

func TestCreate_Integration_Discount(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	// Set a 10% discount for first 5 invoices.
	env.Sub.DiscountPct = 10
	env.Sub.DiscountQty = 5
	tu.RequireNoErr(t, testDB.Save(&env.Sub))

	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 1000, 1)

	svc := inv.NewService(testDB)
	result, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	tu.RequireNoErr(t, err)

	if result.Invoice.DiscountPct != 10 {
		t.Fatalf("expected 10%% discount, got %d%%", result.Invoice.DiscountPct)
	}
	// 1000 - 10% = 900
	if result.Invoice.Total != 900 {
		t.Fatalf("expected total 900, got %d", result.Invoice.Total)
	}
	if result.Invoice.DiscountAmount != 100 {
		t.Fatalf("expected discount 100, got %d", result.Invoice.DiscountAmount)
	}
}

func TestCreate_Integration_PersistsItems(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 250, 4)

	svc := inv.NewService(testDB)
	result, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	tu.RequireNoErr(t, err)

	var items []models.InvoiceItem
	tu.RequireNoErr(t, testDB.List(&items, "invoice_id = ?", result.Invoice.ID))
	if len(items) != 1 {
		t.Fatalf("expected 1 invoice item, got %d", len(items))
	}
	if items[0].Total != 1000 {
		t.Fatalf("expected item total 1000, got %d", items[0].Total)
	}
}
