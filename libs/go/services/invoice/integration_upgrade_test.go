//go:build integration

package invoice_test

import (
	"os"
	"testing"

	inv "github.com/useportcall/portcall/libs/go/services/invoice"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreateUpgrade_Integration_HappyPath(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	newPlan := tu.SeedPlan(t, testDB, env.AppID, "published")

	svc := inv.NewService(testDB)
	result, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{
		SubscriptionID:  env.Sub.ID,
		PriceDifference: 1500,
		OldPlanID:       env.Plan.ID,
		NewPlanID:       newPlan.ID,
	})
	tu.RequireNoErr(t, err)

	if result.Invoice.ID == 0 {
		t.Fatal("invoice was not persisted")
	}
	if !result.ShouldPay {
		t.Fatal("expected ShouldPay = true for positive diff")
	}
	if result.Invoice.Total != 1500 {
		t.Fatalf("expected total 1500, got %d", result.Invoice.Total)
	}
	if result.Invoice.Status != "issued" {
		t.Fatalf("expected issued, got %s", result.Invoice.Status)
	}
}

func TestCreateUpgrade_Integration_ZeroDiff(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	newPlan := tu.SeedPlan(t, testDB, env.AppID, "published")

	svc := inv.NewService(testDB)
	result, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{
		SubscriptionID:  env.Sub.ID,
		PriceDifference: 0,
		OldPlanID:       env.Plan.ID,
		NewPlanID:       newPlan.ID,
	})
	tu.RequireNoErr(t, err)

	if result.ShouldPay {
		t.Fatal("expected ShouldPay = false for zero diff")
	}
}

func TestCreateUpgrade_Integration_NoURL_Errors(t *testing.T) {
	os.Unsetenv("INVOICE_APP_URL")

	svc := inv.NewService(testDB)
	_, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{SubscriptionID: 1})
	if err == nil {
		t.Fatal("expected error for missing INVOICE_APP_URL")
	}
}

func TestCreateUpgrade_Integration_MissingSub_Errors(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	svc := inv.NewService(testDB)
	_, err := svc.CreateUpgrade(&inv.CreateUpgradeInput{
		SubscriptionID:  999999,
		PriceDifference: 100,
		OldPlanID:       1,
		NewPlanID:       2,
	})
	if err == nil {
		t.Fatal("expected error for nonexistent subscription")
	}
}
