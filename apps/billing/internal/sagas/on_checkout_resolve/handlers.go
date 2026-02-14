package on_checkout_resolve

import (
	"encoding/json"
	"log"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func stripeWebhookHandler(c server.IContext) error {
	var input payment.StripeWebhookInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	result, err := svc.ProcessStripeWebhook(&input)
	if err != nil || !result.Handled {
		return err
	}
	switch result.Action {
	case "resolve_checkout_session":
		log.Printf("Enqueued resolve_checkout_session for session %s", result.SessionID)
		return ResolveSession.Enqueue(c.Queue(), map[string]any{
			"external_session_id":        result.SessionID,
			"external_payment_method_id": result.PaymentMethodID,
		})
	case "process_stripe_payment_failure":
		if result.Failure == nil {
			return nil
		}
		return on_payment.ProcessStripePaymentFailure.Enqueue(c.Queue(), result.Failure)
	default:
		return nil
	}
}

func braintreeWebhookHandler(c server.IContext) error {
	var input payment.BraintreeWebhookInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	result, err := svc.ProcessBraintreeWebhook(&input)
	if err != nil || !result.Handled {
		return err
	}
	switch result.Action {
	case "resolve_checkout_session":
		log.Printf("Enqueued resolve_checkout_session for braintree session %s", result.SessionID)
		return ResolveSession.Enqueue(c.Queue(), map[string]any{
			"external_session_id":        result.SessionID,
			"external_payment_method_id": result.PaymentMethodID,
		})
	case "process_stripe_payment_failure":
		if result.Failure == nil {
			return nil
		}
		return on_payment.ProcessStripePaymentFailure.Enqueue(c.Queue(), result.Failure)
	default:
		return nil
	}
}
