//go:build integration

package checkout_session_test

import (
	"errors"
	"testing"

	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_BillingRequired_NoAddress(t *testing.T) {
	appID, _, plan := newEnv(t)
	tu.SeedCompany(t, testDB, appID, 0)

	svc := cs.NewService(testDB, noopCrypto{})
	_, err := svc.Create(&cs.CreateInput{
		AppID:                 appID,
		PlanID:                plan.PublicID,
		UserID:                "irrelevant",
		RequireBillingAddress: true,
	})

	var ve *cs.ValidationError
	if err == nil {
		t.Fatal("expected validation error for missing billing address")
	}
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
}

func TestCreate_BillingRequired_WithAddress(t *testing.T) {
	appID, user, plan := newEnv(t)
	addr := tu.SeedAddress(t, testDB, appID)
	tu.SeedCompany(t, testDB, appID, addr.ID)

	svc := cs.NewService(testDB, noopCrypto{})
	result, err := svc.Create(&cs.CreateInput{
		AppID:                 appID,
		PlanID:                plan.PublicID,
		UserID:                user.PublicID,
		CancelURL:             "https://cancel.test",
		RedirectURL:           "https://redirect.test",
		RequireBillingAddress: true,
	})
	tu.RequireNoErr(t, err)

	if result.Session.CompanyAddressID == nil {
		t.Fatal("expected company address to be set")
	}
	if *result.Session.CompanyAddressID != addr.ID {
		t.Fatalf("expected address ID %d, got %d", addr.ID, *result.Session.CompanyAddressID)
	}
}

func TestCreate_AssignsPaymentCustomerID(t *testing.T) {
	appID, user, plan := newEnv(t)
	svc := cs.NewService(testDB, noopCrypto{})

	result, err := svc.Create(&cs.CreateInput{
		AppID:       appID,
		PlanID:      plan.PublicID,
		UserID:      user.PublicID,
		CancelURL:   "https://cancel.test",
		RedirectURL: "https://redirect.test",
	})
	tu.RequireNoErr(t, err)

	if result.Session.User.PaymentCustomerID == "" {
		t.Fatal("expected payment customer ID to be assigned")
	}
}
