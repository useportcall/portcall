//go:build integration

package checkout_session_test

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

// newEnv seeds the minimal FK graph for checkout session tests.
func newEnv(t *testing.T) (appID uint, user models.User, plan models.Plan) {
	t.Helper()
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	user = tu.SeedUser(t, testDB, app.ID)
	plan = tu.SeedPlan(t, testDB, app.ID, "published")
	conn := tu.SeedConnection(t, testDB, app.ID)
	tu.SeedAppConfig(t, testDB, app.ID, conn.ID)
	return app.ID, user, plan
}
