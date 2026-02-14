//go:build integration

package payment_link_test

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	tu "github.com/useportcall/portcall/libs/go/services/testutil"
)

var testDB dbx.IORM

type noopCrypto struct{}

var _ cryptox.ICrypto = (*noopCrypto)(nil)

func (noopCrypto) Encrypt(data string) (string, error)            { return data, nil }
func (noopCrypto) Decrypt(data string) (string, error)            { return data, nil }
func (noopCrypto) CompareHash(hashed, plain string) (bool, error) { return hashed == plain, nil }

func TestMain(m *testing.M) {
	res := tu.SetupPostgres()
	defer res.Cleanup()
	testDB = res.DB
	os.Exit(m.Run())
}

func newEnv(t *testing.T) (uint, models.User, models.Plan) {
	t.Helper()
	acct := tu.SeedAccount(t, testDB)
	app := tu.SeedApp(t, testDB, acct.ID)
	user := tu.SeedUser(t, testDB, app.ID)
	plan := tu.SeedPlan(t, testDB, app.ID, "published")
	conn := tu.SeedConnection(t, testDB, app.ID)
	tu.SeedAppConfig(t, testDB, app.ID, conn.ID)
	return app.ID, user, plan
}
