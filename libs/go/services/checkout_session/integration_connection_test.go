//go:build integration

package checkout_session_test

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

func TestCreate_FallbackPrefersLocalConnection(t *testing.T) {
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	user := tu.SeedUser(t, testDB, app.ID)
	plan := tu.SeedPlan(t, testDB, app.ID, "published")

	bad := models.Connection{
		PublicID: dbx.GenPublicID("conn"), AppID: app.ID, Name: "Bad",
		Source: "unsupported", PublicKey: "pk_bad", EncryptedKey: "bad",
	}
	if err := testDB.Create(&bad); err != nil {
		t.Fatalf("seed bad connection: %v", err)
	}
	local := tu.SeedConnection(t, testDB, app.ID)
	tu.SeedAppConfig(t, testDB, app.ID, local.ID)
	if err := testDB.Exec(
		"UPDATE app_configs SET default_connection_id = NULL WHERE app_id = ?",
		app.ID,
	); err != nil {
		t.Fatalf("clear default connection: %v", err)
	}

	svc := cs.NewService(testDB, noopCrypto{})
	result, err := svc.Create(&cs.CreateInput{
		AppID: app.ID, PlanID: plan.PublicID, UserID: user.PublicID,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Session.ExternalProvider != local.Source {
		t.Fatalf("expected local provider, got %s", result.Session.ExternalProvider)
	}
}
