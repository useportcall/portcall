//go:build integration

package invoice_test

import (
	"os"
	"testing"

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

type invoiceEnv struct {
	AppID uint
	User  models.User
	Plan  models.Plan
	Addr  models.Address
	Sub   models.Subscription
}

func newEnv(t *testing.T) invoiceEnv {
	t.Helper()
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	user := tu.SeedUser(t, testDB, app.ID)
	plan := tu.SeedPlan(t, testDB, app.ID, "published")
	addr := tu.SeedAddress(t, testDB, app.ID)
	tu.SeedCompany(t, testDB, app.ID, addr.ID)
	planID := plan.ID
	sub := tu.SeedSubscription(t, testDB, app.ID, user.ID, &addr.ID, &planID)
	return invoiceEnv{AppID: app.ID, User: user, Plan: plan, Addr: addr, Sub: sub}
}
