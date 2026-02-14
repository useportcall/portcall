//go:build integration

package checkout_session_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func insertSession(
	t *testing.T, db dbx.IORM, appID, userID, planID uint, status string,
) models.CheckoutSession {
	t.Helper()
	cancel, redirect := "https://cancel", "https://redirect"
	s := models.CheckoutSession{
		PublicID:             dbx.GenPublicID("cs"),
		AppID:                appID,
		UserID:               userID,
		PlanID:               planID,
		ExpiresAt:            time.Now().Add(24 * time.Hour),
		ExternalSessionID:    dbx.GenPublicID("si"),
		ExternalClientSecret: "secret",
		ExternalPublicKey:    "pk",
		ExternalProvider:     "local",
		Status:               status,
		CancelURL:            &cancel,
		RedirectURL:          &redirect,
	}
	if err := db.Create(&s); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	return s
}

func TestResolve_Integration_Active(t *testing.T) {
	appID, user, plan := newEnv(t)
	session := insertSession(t, testDB, appID, user.ID, plan.ID, "active")
	svc := cs.NewService(testDB, noopCrypto{})

	result, err := svc.Resolve(&cs.ResolvePayload{
		ExternalSessionID:       session.ExternalSessionID,
		ExternalPaymentMethodID: "pm_abc",
	})
	tu.RequireNoErr(t, err)

	if result.Skipped {
		t.Fatal("expected resolved, got skipped")
	}
	if result.Session.Status != "resolved" {
		t.Fatalf("expected resolved, got %s", result.Session.Status)
	}

	var reloaded models.CheckoutSession
	tu.RequireNoErr(t, testDB.FindForID(session.ID, &reloaded))
	if reloaded.Status != "resolved" {
		t.Fatalf("DB status expected resolved, got %s", reloaded.Status)
	}
}
