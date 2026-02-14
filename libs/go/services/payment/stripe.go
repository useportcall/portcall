package payment

import (
	"log"

	"github.com/stripe/stripe-go"
)

// ProcessStripeWebhook handles Stripe webhook events.
// No direct mutations - routes to appropriate handlers.
func (s *service) ProcessStripeWebhook(input *StripeWebhookInput) (*StripeResult, error) {
	event := input.Event

	switch event.Type {
	case "setup_intent.succeeded":
		var data stripe.SetupIntent
		if err := data.UnmarshalJSON(event.Data.Raw); err != nil {
			return nil, err
		}

		return &StripeResult{
			Action:          "resolve_checkout_session",
			SessionID:       data.ID,
			PaymentMethodID: data.PaymentMethod.ID,
			Handled:         true,
		}, nil
	case "payment_intent.payment_failed",
		"payment_intent.canceled",
		"payment_intent.requires_action",
		"charge.failed",
		"invoice.payment_failed",
		"invoice.payment_action_required":
		return processStripeFailureEvent(event)
	default:
		log.Println("UNHANDLED_STRIPE_WEBHOOK_TYPE:", event.Type)
		return &StripeResult{Handled: false}, nil
	}
}
