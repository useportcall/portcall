package checkout_session

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// CompleteCheckoutSessionRequest is the request body for completing
// a checkout session after a successful client-side payment.
type CompleteCheckoutSessionRequest struct {
	PaymentMethodID string `json:"payment_method_id"`
}

// CompleteCheckoutSession handles client-side checkout completion.
// For Stripe, subscription creation is webhook-driven after
// setup_intent.succeeded verification. For local/mock providers,
// this endpoint enqueues resolve_checkout_session directly.
func CompleteCheckoutSession(c *routerx.Context, session *models.CheckoutSession) {

	var body CompleteCheckoutSessionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request body")
		return
	}

	// Stripe completion is webhook-driven to ensure subscription activation
	// happens only after a signed setup_intent.succeeded event.
	if session.ExternalProvider == "stripe" {
		c.OK(map[string]any{"status": "processing"})
		return
	}

	paymentMethodID := body.PaymentMethodID
	if paymentMethodID == "" {
		paymentMethodID = "pm_mock_" + dbx.GenPublicID("mock")
	}

	payload := map[string]any{
		"external_session_id":        session.ExternalSessionID,
		"external_payment_method_id": paymentMethodID,
	}

	if err := c.Queue().Enqueue(
		"resolve_checkout_session", payload, "billing_queue",
	); err != nil {
		log.Printf("[CompleteCheckout] enqueue failed: %v", err)
		c.ServerError("failed to complete checkout", err)
		return
	}

	c.OK(map[string]any{"status": "processing"})
}
