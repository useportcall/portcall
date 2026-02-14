package on_subscription_create_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_subscription_create.Steps,
		"create_subscription",
		"create_subscription_items",
		"create_entitlements",
		"create_single_entitlement",
	); err != nil {
		t.Fatal(err)
	}
}
