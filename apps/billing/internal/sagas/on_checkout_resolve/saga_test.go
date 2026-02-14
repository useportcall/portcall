package on_checkout_resolve_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_checkout_resolve.Steps,
		"process_stripe_webhook_event",
		"process_braintree_webhook_event",
		"resolve_checkout_session",
		"create_payment_method",
		"upsert_subscription",
	); err != nil {
		t.Fatal(err)
	}
}
