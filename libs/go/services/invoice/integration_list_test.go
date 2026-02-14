//go:build integration

package invoice_test

import (
	"os"
	"testing"

	inv "github.com/useportcall/portcall/libs/go/services/invoice"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestList_Integration_AllForApp(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 100, 1)

	svc := inv.NewService(testDB)
	tu.RequireNoErr(t, mustCreate(t, svc, env.Sub.ID))

	result, err := svc.List(&inv.ListInput{AppID: env.AppID})
	tu.RequireNoErr(t, err)

	if len(result.Invoices) == 0 {
		t.Fatal("expected at least one invoice")
	}
}

func TestList_Integration_ByUser(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 100, 1)

	svc := inv.NewService(testDB)
	tu.RequireNoErr(t, mustCreate(t, svc, env.Sub.ID))

	result, err := svc.List(&inv.ListInput{
		AppID:  env.AppID,
		UserID: env.User.PublicID,
	})
	tu.RequireNoErr(t, err)

	if len(result.Invoices) == 0 {
		t.Fatal("expected invoices for user")
	}
}

func TestList_Integration_BySub(t *testing.T) {
	os.Setenv("INVOICE_APP_URL", "https://app.test")
	defer os.Unsetenv("INVOICE_APP_URL")

	env := newEnv(t)
	tu.SeedSubscriptionItem(t, testDB, env.AppID, env.Sub.ID, "fixed", 100, 1)

	svc := inv.NewService(testDB)
	tu.RequireNoErr(t, mustCreate(t, svc, env.Sub.ID))

	result, err := svc.List(&inv.ListInput{
		AppID:          env.AppID,
		SubscriptionID: env.Sub.PublicID,
	})
	tu.RequireNoErr(t, err)

	if len(result.Invoices) == 0 {
		t.Fatal("expected invoices for subscription")
	}
}

func TestList_Integration_MissingUser_Errors(t *testing.T) {
	env := newEnv(t)
	svc := inv.NewService(testDB)

	_, err := svc.List(&inv.ListInput{
		AppID:  env.AppID,
		UserID: "usr_nope",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

// mustCreate creates an invoice and returns  only the error.
func mustCreate(t *testing.T, svc inv.Service, subID uint) error {
	t.Helper()
	_, err := svc.Create(&inv.CreateInput{SubscriptionID: subID})
	return err
}
