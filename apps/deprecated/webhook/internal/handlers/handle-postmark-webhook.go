package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// HandlePostmarkWebhook handles incoming webhooks from Postmark.
// It verifies the connection and enqueues the event for processing.
func HandlePostmarkWebhook(c *routerx.Context) {
	log.Printf("[HandlePostmarkWebhook] Received webhook request")

	const MaxBodyBytes = int64(1 * 1024 * 1024) // 1MB limit
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[HandlePostmarkWebhook] Error reading request body: %v", err)
		respondGenericError(c, http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	connectionID, ok := c.Params.Get("connection_id")
	if !ok {
		log.Printf("[HandlePostmarkWebhook] Missing connection_id parameter")
		respondGenericError(c, http.StatusBadRequest)
		return
	}

	log.Printf("[HandlePostmarkWebhook] Processing webhook for connection: %s", connectionID)

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		log.Printf("[HandlePostmarkWebhook] Connection not found: %s, error: %v", connectionID, err)
		respondGenericError(c, http.StatusBadRequest)
		return
	}

	// Basic validation that payload is JSON
	if !json.Valid(payload) {
		log.Printf("[HandlePostmarkWebhook] Invalid JSON payload")
		respondGenericError(c, http.StatusBadRequest)
		return
	}

	// We enqueue the raw JSON payload. The worker will handle unmarshaling
	// into specific event types (Bounce, SpamComplaint, SubscriptionChange, etc.)
	// and verifying any secrets if necessary (Postmark doesn't sign payloads like Stripe
	// by default in the same way, but we rely on the unique connection_id URL).

	jobPayload := map[string]any{
		"connection_id": connection.ID,
		"raw_event":     json.RawMessage(payload),
	}

	if err := c.Queue().Enqueue("process_postmark_webhook_event", jobPayload, "email_queue"); err != nil {
		log.Printf("[HandlePostmarkWebhook] Error enqueueing event: %v", err)
		respondGenericError(c, http.StatusInternalServerError)
		return
	}

	log.Printf("[HandlePostmarkWebhook] Successfully enqueued event for processing")
	c.OK(map[string]any{"status": "success"})
}
