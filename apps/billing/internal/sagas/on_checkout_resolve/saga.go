package on_checkout_resolve

import (
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_create"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_subscription_update"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// --- Flow: process_(stripe|braintree)_webhook_event
//           ├─ resolve_checkout_session → create_payment_method → upsert_subscription
//             ├─ on_subscription_create.CreateSubscription
//             └─ on_subscription_update.UpdateSubscription ---
//           └─ on_payment.ProcessStripePaymentFailure ---

var (
	StripeWebhook       = saga.Route{Name: "process_stripe_webhook_event", Queue: saga.BillingQueue}
	BraintreeWebhook    = saga.Route{Name: "process_braintree_webhook_event", Queue: saga.BillingQueue}
	ResolveSession      = saga.Route{Name: "resolve_checkout_session", Queue: saga.BillingQueue}
	CreatePaymentMethod = saga.Route{Name: "create_payment_method", Queue: saga.BillingQueue}
	UpsertSubscription  = saga.Route{Name: "upsert_subscription", Queue: saga.BillingQueue}
)

var Steps = []saga.Step{
	{Route: StripeWebhook, Handler: stripeWebhookHandler, Emits: []saga.Route{
		ResolveSession,
		on_payment.ProcessStripePaymentFailure,
	}},
	{Route: BraintreeWebhook, Handler: braintreeWebhookHandler, Emits: []saga.Route{
		ResolveSession,
		on_payment.ProcessStripePaymentFailure,
	}},
	{Route: ResolveSession, Handler: resolveSessionHandler, Emits: []saga.Route{CreatePaymentMethod}},
	{Route: CreatePaymentMethod, Handler: createPaymentMethodHandler, Emits: []saga.Route{UpsertSubscription}},
	{Route: UpsertSubscription, Handler: upsertHandler, Emits: []saga.Route{
		on_subscription_create.CreateSubscription,
		on_subscription_update.UpdateSubscription,
	}},
}

func Register(srv server.IServer) { saga.Register(srv, Steps) }
