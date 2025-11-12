package main

import (
	"github.com/useportcall/portcall/apps/billing/handlers"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func main() {
	envx.Load()

	db := dbx.New()
	crypto := cryptox.New()

	server := server.New(db, crypto, map[string]int{
		"billing_queue": 10,
	})

	server.H("create_payment_method", handlers.CreatePaymentMethod)
	server.H("create_subscription", handlers.CreateSubscription)
	server.H("create_subscription_items", handlers.CreateSubscriptionItems)
	server.H("create_invoice", handlers.CreateInvoice)
	server.H("create_invoice_items", handlers.CreateInvoiceItems)
	server.H("calculate_invoice_totals", handlers.CalculateInvoiceTotals)
	server.H("pay_invoice", handlers.PayInvoice)
	server.H("resolve_invoice", handlers.ResolveInvoice)
	server.H("find_subscriptions_to_reset", handlers.FindSubscriptionsToReset)
	server.H("start_subscription_reset", handlers.StartSubscriptionReset)
	server.H("end_subscription_reset", handlers.EndSubscriptionReset)
	server.H("process_meter_event", handlers.ProcessMeterEvent)
	server.H("create_entitlements", handlers.CreateEntitlements)

	server.R()
}
