package webhookx

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func (w *Router) HandlePostmarkWebhook(c *routerx.Context) {
	log.Printf("[webhookx] postmark webhook request")
	if !w.guard(c, "postmark") {
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1*1024*1024)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[webhookx] postmark read body: %v", err)
		respondGenericError(c, http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	connectionID, ok := c.Params.Get("connection_id")
	if !ok || !json.Valid(payload) {
		log.Printf("[webhookx] postmark invalid request")
		respondGenericError(c, http.StatusBadRequest)
		return
	}

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "public_id = ?", connectionID); err != nil {
		log.Printf("[webhookx] postmark missing connection: %s", connectionID)
		respondGenericError(c, http.StatusBadRequest)
		return
	}

	jobPayload := map[string]any{
		"connection_id": connection.ID,
		"raw_event":     json.RawMessage(payload),
	}
	if err := c.Queue().Enqueue("process_postmark_webhook_event", jobPayload, "email_queue"); err != nil {
		log.Printf("[webhookx] postmark enqueue failed: %v", err)
		respondGenericError(c, http.StatusInternalServerError)
		return
	}
	c.OK(map[string]any{"status": "success"})
}
