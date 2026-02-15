package e2etest

import (
	"testing"

	billing "github.com/useportcall/portcall/apps/billing/app"
	"github.com/useportcall/portcall/libs/go/services/checkout_session"
)

// TestE2E_CheckoutResolve_NewCustomer runs the full checkout-resolve saga
// against a real temporary Postgres database. The database is created
// before the test and dropped when it finishes.
func TestE2E_CheckoutResolve_NewCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	t.Setenv("INVOICE_APP_URL", "https://e2e.example.com")

	db := NewTestDB(t)
	f := SeedCheckout(t, db)

	runner := billing.NewFullRunner(db, nil)

	err := runner.Run("resolve_checkout_session", checkout_session.ResolvePayload{
		ExternalSessionID:       f.Session.ExternalSessionID,
		ExternalPaymentMethodID: "pm_e2e_001",
	})
	if err != nil {
		t.Fatalf("saga chain failed: %v", err)
	}

	expect := []string{
		"resolve_checkout_session",
		"create_payment_method",
		"upsert_subscription",
		"create_subscription",
		"pay_invoice",
	}
	for _, name := range expect {
		if !runner.HasTask(name) {
			t.Errorf("missing task %q in chain: %v", name, runner.Executed)
		}
	}

	if runner.HasTask("update_subscription") {
		t.Error("should create, not update, for a new customer")
	}
}
