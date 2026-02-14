//go:build integration

package invoice_test

import (
	"os"
	"testing"

	inv "github.com/useportcall/portcall/libs/go/services/invoice"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_Integration_HappyPath(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 1000, 1)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 500, 2)

	svc := inv.NewService(testDB)
	result, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	tu.RequireNoErr(t, err)

	if result.Skipped {
		t.Fatal("expected created, not skipped")
	}
	if result.Invoice.ID == 0 {
		t.Fatal("invoice was not persisted")
	}
	if result.Invoice.Status != "issued" {
		t.Fatalf("expected issued, got %s", result.Invoice.Status)
	}
	// 1*1000 + 2*500 = 2000
	if result.Invoice.Total != 2000 {
		t.Fatalf("expected total 2000, got %d", result.Invoice.Total)
	}
	if result.Invoice.InvoiceNumber != "INV-0000001" {
		t.Fatalf("expected INV-0000001, got %s", result.Invoice.InvoiceNumber)
	}
}

func TestCreate_Integration_Idempotent(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 100, 1)

	svc := inv.NewService(testDB)

	r1, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	tu.RequireNoErr(t, err)
	if r1.Skipped {
		t.Fatal("first call should create")
	}

	r2, err := svc.Create(&inv.CreateInput{SubscriptionID: env.Sub.ID})
	tu.RequireNoErr(t, err)
	if !r2.Skipped {
		t.Fatal("second call should skip (idempotent)")
	}
}
