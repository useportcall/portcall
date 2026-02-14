//go:build integration

package entitlement_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	ent "github.com/useportcall/portcall/libs/go/services/entitlement"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestResetAll_Integration_HappyPath(t *testing.T) {
	env := newEnv(t)
	e1 := seedEntitlement(t, env.AppID, env.User.ID, 100)
	e2 := seedEntitlement(t, env.AppID, env.User.ID, 50)

	svc := ent.NewService(testDB)
	result, err := svc.ResetAll(&ent.ResetAllInput{UserID: env.User.ID})
	tu.RequireNoErr(t, err)

	if result.ResetCount != 2 {
		t.Fatalf("expected 2 resets, got %d", result.ResetCount)
	}

	// Verify entitlements were actually reset
	var updated1, updated2 models.Entitlement
	testDB.FindForID(e1.ID, &updated1)
	testDB.FindForID(e2.ID, &updated2)

	if updated1.Usage != 0 {
		t.Fatalf("e1 usage should be 0, got %d", updated1.Usage)
	}
	if updated2.Usage != 0 {
		t.Fatalf("e2 usage should be 0, got %d", updated2.Usage)
	}
}

func TestResetAll_Integration_NoEntitlements(t *testing.T) {
	env := newEnv(t)

	svc := ent.NewService(testDB)
	result, err := svc.ResetAll(&ent.ResetAllInput{UserID: env.User.ID})
	tu.RequireNoErr(t, err)

	if result.ResetCount != 0 {
		t.Fatalf("expected 0 resets, got %d", result.ResetCount)
	}
}

func TestResetAll_Integration_OnlyResetsTargetUser(t *testing.T) {
	env := newEnv(t)
	otherUser := tu.SeedUser(t, testDB, env.AppID)
	
	seedEntitlement(t, env.AppID, env.User.ID, 100)
	otherEnt := seedEntitlement(t, env.AppID, otherUser.ID, 200)

	svc := ent.NewService(testDB)
	result, err := svc.ResetAll(&ent.ResetAllInput{UserID: env.User.ID})
	tu.RequireNoErr(t, err)

	if result.ResetCount != 1 {
		t.Fatalf("expected 1 reset, got %d", result.ResetCount)
	}

	// Other user's entitlement should be unchanged
	var unchanged models.Entitlement
	testDB.FindForID(otherEnt.ID, &unchanged)
	if unchanged.Usage != 200 {
		t.Fatalf("other user's ent should be 200, got %d", unchanged.Usage)
	}
}
