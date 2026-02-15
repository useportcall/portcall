package harness

import (
	"testing"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// SeedBase creates the account, app (test mode), address, company,
// local connection, and app config required as prerequisites.
func (h *Harness) SeedBase() {
	h.T.Helper()

	acct := &models.Account{Email: "admin@portcall.internal", FirstName: "E2E", LastName: "Test"}
	Must(h.T, h.DB.Create(acct))

	app := &models.App{
		PublicID: "app_e2e_test", Name: "E2E App",
		Status: "active", AccountID: acct.ID,
		PublicApiKey: "pk_e2e_test", BillingExempt: true,
	}
	Must(h.T, h.DB.Create(app))

	addr := &models.Address{
		PublicID: "addr_e2e_test", AppID: app.ID,
		Line1: "1 Test Blvd", City: "San Francisco",
		PostalCode: "94105", Country: "US",
	}
	Must(h.T, h.DB.Create(addr))

	company := &models.Company{
		AppID: app.ID, Name: "E2E Corp",
		BillingAddressID: addr.ID, Email: "billing@e2e.test",
	}
	Must(h.T, h.DB.Create(company))

	conn := &models.Connection{
		PublicID: "connect_e2e_test", AppID: app.ID,
		Source: "local", Name: "Local Dev", PublicKey: "pk_test",
	}
	Must(h.T, h.DB.Create(conn))

	cfg := &models.AppConfig{AppID: app.ID, DefaultConnectionID: conn.ID}
	Must(h.T, h.DB.Create(cfg))

	h.AppPublicID = app.PublicID
	h.AppID = app.ID
}

// SeedCheckout inserts the minimal dataset for checkout-resolve saga tests.
func SeedCheckout(t *testing.T, db dbx.IORM) *CheckoutFixture {
	t.Helper()
	f := &CheckoutFixture{}

	f.Account = &models.Account{Email: "e2e@test.com", FirstName: "E2E", LastName: "Test"}
	Must(t, db.Create(f.Account))

	f.App = &models.App{PublicID: "app_e2e", Name: "E2E App", Status: "active", AccountID: f.Account.ID}
	Must(t, db.Create(f.App))

	f.Address = &models.Address{
		PublicID: "addr_e2e", AppID: f.App.ID, Line1: "123 E2E St",
		City: "Testville", PostalCode: "12345", Country: "US",
	}
	Must(t, db.Create(f.Address))

	f.Company = &models.Company{AppID: f.App.ID, Name: "E2E Co", BillingAddressID: f.Address.ID}
	Must(t, db.Create(f.Company))

	f.User = &models.User{
		PublicID: "usr_e2e", AppID: f.App.ID, Name: "Jane",
		Email: "jane@e2e.test", PaymentCustomerID: "cus_e2e",
		BillingAddressID: &f.Address.ID,
	}
	Must(t, db.Create(f.User))

	f.Connection = &models.Connection{AppID: f.App.ID, Source: "local", PublicKey: "pk_e2e"}
	Must(t, db.Create(f.Connection))

	f.Plan = &models.Plan{
		PublicID: "plan_e2e", AppID: f.App.ID, Name: "Pro",
		Status: "published", Interval: "month", IntervalCount: 1,
		Currency: "usd", InvoiceDueByDays: 14,
	}
	Must(t, db.Create(f.Plan))

	f.PlanItem = &models.PlanItem{
		PublicID: "pi_e2e", AppID: f.App.ID, PlanID: f.Plan.ID,
		PricingModel: "fixed", UnitAmount: 2999, Quantity: 1,
		PublicTitle: "Pro Plan",
	}
	Must(t, db.Create(f.PlanItem))

	f.Feature = &models.Feature{PublicID: "feat_e2e", AppID: f.App.ID, IsMetered: false}
	Must(t, db.Create(f.Feature))

	f.PlanFeature = &models.PlanFeature{
		PublicID: "pf_e2e", AppID: f.App.ID, PlanID: f.Plan.ID,
		PlanItemID: f.PlanItem.ID, FeatureID: f.Feature.ID,
		Interval: "month", Quota: 10000,
	}
	Must(t, db.Create(f.PlanFeature))

	redirect := "https://e2e.test/ok"
	cancel := "https://e2e.test/cancel"
	f.Session = &models.CheckoutSession{
		PublicID: "cs_e2e", AppID: f.App.ID, PlanID: f.Plan.ID,
		UserID: f.User.ID, ExternalSessionID: "seti_e2e",
		ExternalClientSecret: "secret", ExternalPublicKey: "pk",
		ExternalProvider: "local", ExpiresAt: time.Now().Add(time.Hour),
		RedirectURL: &redirect,
		CancelURL:   &cancel, Status: "active",
	}
	Must(t, db.Create(f.Session))
	return f
}

// CheckoutFixture holds all seeded entities for checkout-resolve tests.
type CheckoutFixture struct {
	Account     *models.Account
	App         *models.App
	Address     *models.Address
	Company     *models.Company
	User        *models.User
	Connection  *models.Connection
	Plan        *models.Plan
	PlanItem    *models.PlanItem
	Feature     *models.Feature
	PlanFeature *models.PlanFeature
	Session     *models.CheckoutSession
}

// Must fails the test if err is non-nil.
func Must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
}
