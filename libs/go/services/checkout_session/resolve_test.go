package checkout_session_test

import (
	"errors"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
)

// resolveDB is a minimal dbx.IORM mock for resolve tests.
type resolveDB struct {
	dbx.IORM // embedded to satisfy interface; only FindFirst and Save are used
	session  *models.CheckoutSession
	findErr  error
	saveErr  error
	saved    string
}

func (m *resolveDB) FindFirst(dest any, conds ...any) error {
	if m.findErr != nil {
		return m.findErr
	}
	s := dest.(*models.CheckoutSession)
	*s = *m.session
	return nil
}

func (m *resolveDB) Save(dest any) error {
	s := dest.(*models.CheckoutSession)
	m.saved = s.Status
	return m.saveErr
}

func TestResolve_ActiveSession(t *testing.T) {
	db := &resolveDB{
		session: &models.CheckoutSession{Status: "active", AppID: 1, UserID: 2, PlanID: 10},
	}
	svc := cs.NewService(db, nil)

	result, err := svc.Resolve(&cs.ResolvePayload{
		ExternalSessionID:       "si_123",
		ExternalPaymentMethodID: "pm_456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Skipped {
		t.Fatal("expected resolved, got skipped")
	}
	if result.ExternalPaymentMethodID != "pm_456" {
		t.Fatalf("expected pm_456, got %s", result.ExternalPaymentMethodID)
	}
	if db.saved != "resolved" {
		t.Fatalf("expected status resolved, got %s", db.saved)
	}
}

func TestResolve_InactiveSession_Skips(t *testing.T) {
	db := &resolveDB{
		session: &models.CheckoutSession{Status: "completed"},
	}
	svc := cs.NewService(db, nil)

	result, err := svc.Resolve(&cs.ResolvePayload{ExternalSessionID: "si_123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Skipped {
		t.Fatal("expected skipped for non-active session")
	}
}

func TestResolve_FindError(t *testing.T) {
	db := &resolveDB{
		findErr: errors.New("not found"),
	}
	svc := cs.NewService(db, nil)

	_, err := svc.Resolve(&cs.ResolvePayload{ExternalSessionID: "si_bad"})
	if err == nil {
		t.Fatal("expected error")
	}
}
