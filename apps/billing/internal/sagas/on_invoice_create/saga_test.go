package on_invoice_create_test

import (
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestFlow(t *testing.T) {
	if err := saga.ExpectRoutes(on_invoice_create.Steps,
		"create_invoice",
	); err != nil {
		t.Fatal(err)
	}
}
