package on_entitlement_upsert_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_entitlement_upsert"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_entitlement_upsert.Steps,
		"start_entitlement_upsert",
		"upsert_entitlement",
	); err != nil {
		t.Fatal(err)
	}
}
