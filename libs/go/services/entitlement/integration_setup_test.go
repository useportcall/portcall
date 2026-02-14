//go:build integration

package entitlement_test

import (
	"os"
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

var testDB dbx.IORM

func TestMain(m *testing.M) {
	res := tu.SetupPostgres()
	defer res.Cleanup()
	testDB = res.DB
	os.Exit(m.Run())
}

type entEnv struct {
	AppID uint
	User  models.User
}

func newEnv(t *testing.T) entEnv {
	t.Helper()
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	user := tu.SeedUser(t, testDB, app.ID)
	return entEnv{AppID: app.ID, User: user}
}

func seedEntitlement(t *testing.T, appID, userID uint, usage int64) models.Entitlement {
	t.Helper()
	// Create a feature first (required for FeaturePublicID foreign key)
	feature := models.Feature{
		PublicID: dbx.GenPublicID("feat"),
		AppID:    appID,
	}
	if err := testDB.Create(&feature); err != nil {
		t.Fatalf("seed feature: %v", err)
	}

	now := time.Now()
	e := models.Entitlement{
		AppID:           appID,
		UserID:          userID,
		FeaturePublicID: feature.PublicID,
		Interval:        "monthly",
		Quota:           1000,
		Usage:           usage,
		LastResetAt:     &now,
		AnchorAt:        &now,
	}
	if err := testDB.Create(&e); err != nil {
		t.Fatalf("seed entitlement: %v", err)
	}
	return e
}
