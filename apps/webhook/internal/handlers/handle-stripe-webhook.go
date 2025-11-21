package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/webhook"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func HandleStripeWebhook(c *routerx.Context) {
	const MaxBodyBytes = int64(65536) // Limit the body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.ServerError("READ_REQUEST_BODY_ERROR", err)
		return
	}
	defer c.Request.Body.Close()

	connectionID, ok := c.Params.Get("connection_id")
	if !ok {
		c.ServerError("NO_CONNECTION_ID_IN_PARAMS", fmt.Errorf("no connection_id in params"))
		return
	}

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		c.ServerError("NO_STRIPE_CONNECTION_FOR_CONNECTION", err)
		return
	}

	if connection.EncryptedWebhookSecret == nil {
		c.ServerError("NO_STRIPE_WEBHOOK_SECRET_FOR_CONNECTION", fmt.Errorf("no stripe webhook secret for connection"))
		return
	}

	secret, err := c.Crypto().Decrypt(*connection.EncryptedWebhookSecret)
	if err != nil {
		c.ServerError("DECRYPT_WEBHOOK_SECRET_ERROR", err)
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sigHeader, secret)
	if err != nil {
		c.ServerError("INVALID_WEBHOOK_SIGNATURE", err)
		return
	}

	if err := c.Queue().Enqueue("process_stripe_webhook_event", event, "billing_queue"); err != nil {
		c.ServerError("ENQUEUE_STRIPE_WEBHOOK_EVENT_ERROR", err)
		return
	}

	c.OK(map[string]any{"status": "success"})
}
