package handlers

import (
	"encoding/json"
	"log"

	"github.com/stripe/stripe-go"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func ProcessStripeWebhookEvent(c server.IContext) error {
	var event stripe.Event
	if err := json.Unmarshal(c.Payload(), &event); err != nil {
		return err
	}

	switch event.Type {
	case "setup_intent.succeeded":
		var data stripe.SetupIntent
		if err := json.Unmarshal(event.Data.Raw, &data); err != nil {
			return err
		}

		payload := CreatePaymentMethodPayload{
			ExternalSessionID:       data.ID,
			ExternalPaymentMethodID: data.PaymentMethod.ID,
		}
		if err := c.Queue().Enqueue("create_payment_method", payload, "billing_queue"); err != nil {
			return err
		}

	default:
		log.Println("UNHANDLED_STRIPE_WEBHOOK_TYPE:", event.Type)
	}

	return nil
}
