package on_entitlement_reset_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_entitlement_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_entitlement_reset.Steps,
		"reset_all_entitlements",
	); err != nil {
		t.Fatal(err)
	}
}
