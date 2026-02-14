//go:build integration

package checkout_session_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_Roundtrip_FullFlow(t *testing.T) {
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

	resolved, err := svc.Resolve(&cs.ResolvePayload{
		ExternalSessionID:       result.Session.ExternalSessionID,
		ExternalPaymentMethodID: "pm_roundtrip",
	})
	tu.RequireNoErr(t, err)

	if resolved.Skipped {
		t.Fatal("expected resolved")
	}
	if resolved.Session.Status != "resolved" {
		t.Fatalf("expected resolved, got %s", resolved.Session.Status)
	}

	var found models.CheckoutSession
	tu.RequireNoErr(t, testDB.FindFirst(&found,
		"public_id = ?", result.Session.PublicID))
	if found.Status != "resolved" {
		t.Fatalf("expected resolved in DB, got %s", found.Status)
	}
}

func TestCreate_InvalidPlanID_Errors(t *testing.T) {
	appID, user, _ := newEnv(t)
	svc := cs.NewService(testDB, noopCrypto{})

	_, err := svc.Create(&cs.CreateInput{
		AppID:       appID,
		PlanID:      "plan_doesnotexist",
		UserID:      user.PublicID,
		CancelURL:   "https://cancel.test",
		RedirectURL: "https://redirect.test",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent plan")
	}
}

func TestCreate_InvalidUserID_Errors(t *testing.T) {
	appID, _, plan := newEnv(t)
	svc := cs.NewService(testDB, noopCrypto{})

	_, err := svc.Create(&cs.CreateInput{
		AppID:       appID,
		PlanID:      plan.PublicID,
		UserID:      "usr_doesnotexist",
		CancelURL:   "https://cancel.test",
		RedirectURL: "https://redirect.test",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}
