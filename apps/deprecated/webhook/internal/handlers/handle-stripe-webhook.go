package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/webhook"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// Relevant Stripe events for card collection, storage, and charging.
// Portcall handles all subscription/billing logic internally - we only use
// Stripe for payment method collection and charging.
var relevantStripeEvents = map[string]bool{
	// Card/Payment Method events
	"setup_intent.succeeded":                    true, // Card saved successfully
	"setup_intent.setup_failed":                 true, // Card save failed
	"payment_intent.succeeded":                  true, // Charge successful
	"payment_intent.payment_failed":             true, // Charge failed
	"payment_method.attached":                   true, // Payment method attached to customer
	"payment_method.detached":                   true, // Payment method removed from customer
	"payment_method.updated":                    true, // Payment method updated (e.g., new expiry)
	"payment_method.card_automatically_updated": true, // Card auto-updated by network
}

func HandleStripeWebhook(c *routerx.Context) {
	log.Printf("[HandleStripeWebhook] Received webhook request")

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[HandleStripeWebhook] Error reading request body: %v", err)
		respondOK(c) // Always return OK to Stripe - never leak error info
		return
	}
	defer c.Request.Body.Close()

	connectionID, ok := c.Params.Get("connection_id")
	if !ok {
		log.Printf("[HandleStripeWebhook] Missing connection_id parameter")
		respondOK(c)
		return
	}

	log.Printf("[HandleStripeWebhook] Processing webhook for connection: %s", connectionID)

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		log.Printf("[HandleStripeWebhook] Connection not found: %s, error: %v", connectionID, err)
		respondOK(c)
		return
	}

	log.Printf("[HandleStripeWebhook] Found connection ID=%d, AppID=%d, Source=%s", connection.ID, connection.AppID, connection.Source)

	if connection.EncryptedWebhookSecret == nil {
		log.Printf("[HandleStripeWebhook] Connection %s has no webhook secret configured", connectionID)
		respondOK(c)
		return
	}

	secret, err := c.Crypto().Decrypt(*connection.EncryptedWebhookSecret)
	if err != nil {
		log.Printf("[HandleStripeWebhook] Error decrypting webhook secret: %v", err)
		respondOK(c)
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	if sigHeader == "" {
		log.Printf("[HandleStripeWebhook] Missing Stripe-Signature header")
		respondOK(c)
		return
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, secret)
	if err != nil {
		log.Printf("[HandleStripeWebhook] Error constructing webhook event (signature verification failed): %v", err)
		respondOK(c)
		return
	}

	log.Printf("[HandleStripeWebhook] Received event: ID=%s, Type=%s", event.ID, event.Type)

	// Only process relevant card/payment events - ignore subscription events
	// as Portcall handles all subscription logic internally
	if !relevantStripeEvents[string(event.Type)] {
		log.Printf("[HandleStripeWebhook] Ignoring non-relevant event type: %s", event.Type)
		respondOK(c)
		return
	}

	if err := c.Queue().Enqueue("process_stripe_webhook_event", event, "billing_queue"); err != nil {
		log.Printf("[HandleStripeWebhook] Error enqueueing event: %v", err)
		// Still return OK - we log the error but don't expose it
		respondOK(c)
		return
	}

	log.Printf("[HandleStripeWebhook] Successfully enqueued event %s for processing", event.ID)
	respondOK(c)
}

// respondOK always returns a 200 OK response to Stripe.
// We never return error information to prevent information leakage.
// All errors are logged internally for debugging.
func respondOK(c *routerx.Context) {
	c.JSON(http.StatusOK, map[string]any{"received": true})
}
