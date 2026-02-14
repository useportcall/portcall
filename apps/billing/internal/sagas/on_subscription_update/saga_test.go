package on_subscription_update_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_subscription_update.Steps,
		"update_subscription",
		"process_plan_switch",
		"create_upgrade_invoice",
	); err != nil {
		t.Fatal(err)
	}
}
