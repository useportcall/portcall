package on_subscription_reset

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: find_subscriptions_to_reset → start_subscription_reset
//             ├─ (active)   → on_invoice_create.CreateInvoice + end_subscription_reset
//             └─ (canceled) → on_subscription_create.CreateSubscription (rollback plan)
//                             + subscription_rollback_email ---

var rollbackEmail = saga.Route{Name: "subscription_rollback_email", Queue: saga.EmailQueue}

var (
	FindSubscriptions = saga.Route{Name: "find_subscriptions_to_reset", Queue: saga.BillingQueue}
	StartReset        = saga.Route{Name: "start_subscription_reset", Queue: saga.BillingQueue}
	EndReset          = saga.Route{Name: "end_subscription_reset", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: FindSubscriptions, Handler: findHandler, Emits: []saga.Route{StartReset}},
	{Route: StartReset, Handler: startResetHandler, Emits: []saga.Route{
		on_invoice_create.CreateInvoice, EndReset, // active path
		on_subscription_create.CreateSubscription, rollbackEmail, // canceled path
	}},
	{Route: EndReset, Handler: endResetHandler},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
