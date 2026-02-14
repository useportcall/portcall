package on_subscription_update

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: update_subscription -> [parallel fan-out]
//             |- on_subscription_create.CreateItems (bulk transactional)
//             +- process_plan_switch
//                  +- (upgrade) create_upgrade_invoice -> on_payment.PayInvoice
//                  +- (all changes) on_subscription_create.CreateEntitlements ---

var (
	UpdateSubscription   = saga.Route{Name: "update_subscription", Queue: saga.BillingQueue}
	ProcessPlanSwitch    = saga.Route{Name: "process_plan_switch", Queue: saga.BillingQueue}
	CreateUpgradeInvoice = saga.Route{Name: "create_upgrade_invoice", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: UpdateSubscription, Handler: updateHandler, Emits: []saga.Route{
		on_subscription_create.CreateItems, ProcessPlanSwitch,
	}},
	{Route: ProcessPlanSwitch, Handler: planSwitchHandler, Emits: []saga.Route{
		CreateUpgradeInvoice, on_subscription_create.CreateEntitlements,
	}},
	{Route: CreateUpgradeInvoice, Handler: upgradeInvoiceHandler, Emits: []saga.Route{
		on_payment.PayInvoice,
	}},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
