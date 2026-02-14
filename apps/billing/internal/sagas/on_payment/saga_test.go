package on_payment_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_payment.Steps,
		"pay_invoice",
		"resolve_invoice",
		"process_stripe_payment_failure",
	); err != nil {
		t.Fatal(err)
	}
}
