//go:build integration

package subscription_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

var integrationDB dbx.IORM

func TestMain(m *testing.M) {
	res := tu.SetupPostgres()
	defer res.Cleanup()
	integrationDB = res.DB
	os.Exit(m.Run())
}

func seedIntegrationAppUser(t *testing.T) (models.App, models.User) {
	t.Helper()
	acct := tu.SeedAccount(t, integrationDB)
	app := tu.SeedApp(t, integrationDB, acct.ID)
	user := tu.SeedUser(t, integrationDB, app.ID)
	return app, user
}

func seedPlan(t *testing.T, appID uint, name string, isFree bool) models.Plan {
	t.Helper()
	plan := models.Plan{
		PublicID: dbx.GenPublicID("plan"),
		AppID:    appID,
		Name:     name,
		Status:   "published",
		Interval: "month", IntervalCount: 1,
		Currency: "USD", IsFree: isFree,
	}
	if err := integrationDB.Create(&plan); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	return plan
}

func seedPlanItem(t *testing.T, appID, planID uint, title string, amount int64) {
	t.Helper()
	item := models.PlanItem{
		PublicID: dbx.GenPublicID("pi"),
		AppID:    appID, PlanID: planID,
		PricingModel: "fixed", Quantity: 1, UnitAmount: amount,
		PublicTitle: title, PublicDescription: title, PublicUnitLabel: "unit",
		Interval: "month", IntervalCount: 1,
	}
	if err := integrationDB.Create(&item); err != nil {
		t.Fatalf("seed plan item: %v", err)
	}
}
