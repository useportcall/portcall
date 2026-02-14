package on_payment

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: pay_invoice → resolve_invoice → send_invoice_paid_email (email queue) ---

var sendInvoiceEmail = saga.Route{Name: "send_invoice_paid_email", Queue: saga.EmailQueue}
var sendInvoiceDunningEmail = saga.Route{Name: "send_invoice_dunning_email", Queue: saga.EmailQueue}

var (
	PayInvoice                  = saga.Route{Name: "pay_invoice", Queue: saga.BillingQueue}
	ResolveInvoice              = saga.Route{Name: "resolve_invoice", Queue: saga.BillingQueue}
	ProcessStripePaymentFailure = saga.Route{Name: "process_stripe_payment_failure", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: PayInvoice, Handler: payHandler, Emits: []saga.Route{
		ResolveInvoice,
		ProcessStripePaymentFailure,
	}},
	{Route: ResolveInvoice, Handler: resolveHandler, Emits: []saga.Route{sendInvoiceEmail}},
	{Route: ProcessStripePaymentFailure, Handler: stripePaymentFailureHandler},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
