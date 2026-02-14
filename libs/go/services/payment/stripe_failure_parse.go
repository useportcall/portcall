package payment

import (
	"fmt"

	"github.com/stripe/stripe-go"
)

const stripeInvoiceMetadataKey = "portcall_invoice_id"

func processStripeFailureEvent(event stripe.Event) (*StripeResult, error) {
	payload, err := extractStripeFailurePayload(event)
	if err != nil {
		return nil, err
	}
	if payload == nil {
		return &StripeResult{Handled: true}, nil
	}
	return &StripeResult{
		Action:  "process_stripe_payment_failure",
		Failure: payload,
		Handled: true,
	}, nil
}

func extractStripeFailurePayload(event stripe.Event) (*StripeFailurePayload, error) {
	switch event.Type {
	case "payment_intent.payment_failed", "payment_intent.canceled", "payment_intent.requires_action":
		return parsePaymentIntentFailure(event)
	case "charge.failed":
		return parseChargeFailure(event)
	case "invoice.payment_failed", "invoice.payment_action_required":
		return parseInvoiceFailure(event)
	default:
		return nil, fmt.Errorf("unsupported stripe failure event: %s", event.Type)
	}
}
