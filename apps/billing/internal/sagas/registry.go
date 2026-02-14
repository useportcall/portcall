package sagas

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_checkout_resolve"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_df_decrement"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_df_increment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_entitlement_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_entitlement_upsert"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_meter_event"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_reset"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// All returns every saga's steps for cross-saga validation.
func All() [][]saga.Step {
	return [][]saga.Step{
		on_checkout_resolve.Steps,
		on_subscription_create.Steps,
		on_subscription_update.Steps,
		on_invoice_create.Steps,
		on_payment.Steps,
		on_subscription_reset.Steps,
		on_entitlement_reset.Steps,
		on_entitlement_upsert.Steps,
		on_meter_event.Steps,
		on_df_increment.Steps,
		on_df_decrement.Steps,
	}
}

// RegisterAll registers every saga's handlers with the queue server.
func RegisterAll(srv server.IServer) {
	for _, steps := range All() {
		saga.Register(srv, steps)
	}
}
