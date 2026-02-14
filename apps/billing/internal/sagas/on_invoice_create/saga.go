package on_invoice_create

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: create_invoice (items + totals in txn) → on_payment.PayInvoice ---

var CreateInvoice = saga.Route{Name: "create_invoice", Queue: saga.BillingQueue}

var Steps = []saga.Step{
	{Route: CreateInvoice, Handler: createInvoiceHandler, Emits: []saga.Route{on_payment.PayInvoice}},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
