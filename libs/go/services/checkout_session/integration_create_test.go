//go:build integration

package checkout_session_test

import (
	"errors"
	"testing"

	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_HappyPath(t *testing.T) {
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

	s := result.Session
	if s.ID == 0 {
		t.Fatal("session was not persisted")
	}
	if s.Status != "active" {
		t.Fatalf("expected active, got %s", s.Status)
	}
	if s.ExternalProvider != "local" {
		t.Fatalf("expected local provider, got %s", s.ExternalProvider)
	}
	if s.CompanyAddressID != nil {
		t.Fatalf("expected nil company address, got %v", *s.CompanyAddressID)
	}
}

func TestCreate_DraftPlan_Rejected(t *testing.T) {
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	tu.SeedUser(t, testDB, app.ID)
	draft := tu.SeedPlan(t, testDB, app.ID, "draft")

	svc := cs.NewService(testDB, noopCrypto{})
	_, err := svc.Create(&cs.CreateInput{
		AppID:  app.ID,
		PlanID: draft.PublicID,
		UserID: "nonexistent",
	})

	var ve *cs.ValidationError
	if err == nil {
		t.Fatal("expected validation error for draft plan")
	}
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
}
