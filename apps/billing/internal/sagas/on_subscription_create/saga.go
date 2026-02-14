package on_subscription_create

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_invoice_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: create_subscription (creates subscription + items in one txn)
//             ├─ on_invoice_create.CreateInvoice (only if item_count > 0)
//             └─ create_entitlements → create_single_entitlement (iterative) ---

var (
	CreateSubscription      = saga.Route{Name: "create_subscription", Queue: saga.BillingQueue}
	CreateItems             = saga.Route{Name: "create_subscription_items", Queue: saga.BillingQueue}
	CreateEntitlements      = saga.Route{Name: "create_entitlements", Queue: saga.BillingQueue}
	CreateSingleEntitlement = saga.Route{Name: "create_single_entitlement", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: CreateSubscription, Handler: createSubscriptionHandler, Emits: []saga.Route{CreateItems, on_invoice_create.CreateInvoice, CreateEntitlements}},
	{Route: CreateItems, Handler: createItemsHandler, Emits: []saga.Route{on_invoice_create.CreateInvoice}},
	{Route: CreateEntitlements, Handler: createEntitlementsHandler, Emits: []saga.Route{CreateSingleEntitlement}},
	{Route: CreateSingleEntitlement, Handler: createSingleEntitlementHandler, Emits: []saga.Route{CreateSingleEntitlement}},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
