package entitlement_test

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	ent "github.com/useportcall/portcall/libs/go/services/entitlement"
)

func TestResetAll_Success(t *testing.T) {
	now := time.Now()
	db := &mockDB{
		entitlements: []models.Entitlement{
			{AppID: 1, UserID: 10, Usage: 100, LastResetAt: &now},
			{AppID: 1, UserID: 10, Usage: 50, LastResetAt: &now},
		},
	}
	db.entitlements[0].ID = 1
	db.entitlements[1].ID = 2

	svc := ent.NewService(db)
	result, err := svc.ResetAll(&ent.ResetAllInput{UserID: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ResetCount != 2 {
		t.Fatalf("expected 2 resets, got %d", result.ResetCount)
	}
	if !db.txnCalled {
		t.Fatal("expected transaction to be used")
	}
	if len(db.savedEntitlements) != 2 {
		t.Fatalf("expected 2 saves, got %d", len(db.savedEntitlements))
	}
}

func TestResetAll_NoEntitlements(t *testing.T) {
	db := &mockDB{
		entitlements: []models.Entitlement{},
	}

	svc := ent.NewService(db)
	result, err := svc.ResetAll(&ent.ResetAllInput{UserID: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ResetCount != 0 {
		t.Fatalf("expected 0 resets, got %d", result.ResetCount)
	}
	if db.txnCalled {
		t.Fatal("no transaction needed for empty list")
	}
}
