package webhookx

import (
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/webhook"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func (w *Router) HandleStripeWebhook(c *routerx.Context) {
	log.Printf("[webhookx] stripe webhook request")
	if !w.guard(c, "stripe") {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 65536)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[webhookx] stripe read body: %v", err)
		respondStripeOK(c)
		return
	}
	defer c.Request.Body.Close()

	connectionID, ok := c.Params.Get("connection_id")
	if !ok {
		log.Printf("[webhookx] stripe missing connection_id")
		respondStripeOK(c)
		return
	}

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		log.Printf("[webhookx] stripe missing connection: %s", connectionID)
		respondStripeOK(c)
		return
	}
	if connection.EncryptedWebhookSecret == nil {
		log.Printf("[webhookx] stripe no webhook secret: %s", connectionID)
		respondStripeOK(c)
		return
	}

	secret, err := c.Crypto().Decrypt(*connection.EncryptedWebhookSecret)
	if err != nil || c.GetHeader("Stripe-Signature") == "" {
		log.Printf("[webhookx] stripe invalid secret/signature")
		respondStripeOK(c)
		return
	}
	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), secret)
	if err != nil || !relevantStripeEvents[string(event.Type)] {
		if err != nil {
			log.Printf("[webhookx] stripe signature failed: %v", err)
		}
		respondStripeOK(c)
		return
	}
	if err := c.Queue().Enqueue("process_stripe_webhook_event", event, "billing_queue"); err != nil {
		log.Printf("[webhookx] stripe enqueue failed: %v", err)
		respondStripeError(c)
		return
	}
	respondStripeOK(c)
}
