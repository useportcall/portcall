//go:build integration

package checkout_session_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestResolve_Integration_AlreadyResolved(t *testing.T) {
	appID, user, plan := newEnv(t)
	session := insertSession(t, testDB, appID, user.ID, plan.ID, "resolved")
	svc := cs.NewService(testDB, noopCrypto{})

	result, err := svc.Resolve(&cs.ResolvePayload{
		ExternalSessionID: session.ExternalSessionID,
	})
	tu.RequireNoErr(t, err)

	if !result.Skipped {
		t.Fatal("expected skipped for non-active session")
	}
}

func TestResolve_Integration_NotFound(t *testing.T) {
	svc := cs.NewService(testDB, noopCrypto{})

	_, err := svc.Resolve(&cs.ResolvePayload{
		ExternalSessionID: "nonexistent",
	})
	if err == nil {
		t.Fatal("expected error for missing session")
	}
	if !dbx.IsRecordNotFoundError(err) {
		t.Fatalf("expected record-not-found, got %v", err)
	}
}

func TestResolve_Integration_Idempotent(t *testing.T) {
	appID, user, plan := newEnv(t)
	session := insertSession(t, testDB, appID, user.ID, plan.ID, "active")
	svc := cs.NewService(testDB, noopCrypto{})

	payload := &cs.ResolvePayload{
		ExternalSessionID:       session.ExternalSessionID,
		ExternalPaymentMethodID: "pm_1",
	}

	r1, err := svc.Resolve(payload)
	tu.RequireNoErr(t, err)
	if r1.Skipped {
		t.Fatal("first resolve should not skip")
	}

	r2, err := svc.Resolve(payload)
	tu.RequireNoErr(t, err)
	if !r2.Skipped {
		t.Fatal("second resolve should skip")
	}
}
