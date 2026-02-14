package on_subscription_reset_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_subscription_reset.Steps,
		"find_subscriptions_to_reset",
		"start_subscription_reset",
		"end_subscription_reset",
	); err != nil {
		t.Fatal(err)
	}
}
