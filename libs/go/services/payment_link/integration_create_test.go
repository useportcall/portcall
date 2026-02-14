//go:build integration

package payment_link_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	pl "github.com/useportcall/portcall/libs/go/services/payment_link"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_AutoCreatesUserFromEmail(t *testing.T) {
	appID, _, plan := newEnv(t)
	t.Setenv("CHECKOUT_URL", "https://checkout.test")
	svc := pl.NewService(testDB, noopCrypto{})

	result, err := svc.Create(&pl.CreateInput{
		AppID:       appID,
		PlanID:      plan.PublicID,
		UserEmail:   "new-user@example.com",
		UserName:    "New User",
		CancelURL:   "https://example.com/cancel",
		RedirectURL: "https://example.com/success",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !strings.Contains(result.URL, "pl=") || !strings.Contains(result.URL, "pt=") {
		t.Fatalf("expected signed payment link URL, got %s", result.URL)
	}

	var created models.User
	if err := testDB.FindFirst(&created, "app_id = ? AND id = ?", appID, result.PaymentLink.UserID); err != nil {
		t.Fatalf("FindFirst(user) error = %v", err)
	}
	if created.Email != "new-user@example.com" {
		t.Fatalf("expected created user email to match, got %s", created.Email)
	}
}

func TestCreate_WithoutPaymentConnection_Rejected(t *testing.T) {
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	plan := tu.SeedPlan(t, testDB, app.ID, "published")

	svc := pl.NewService(testDB, noopCrypto{})
	_, err := svc.Create(&pl.CreateInput{
		AppID:     app.ID,
		PlanID:    plan.PublicID,
		UserEmail: "new-user@example.com",
	})

	var ve *pl.ValidationError
	if err == nil {
		t.Fatal("expected validation error for missing payment connection")
	}
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
}
