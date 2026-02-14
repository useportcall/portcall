package connection

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteConnection(c *routerx.Context) {
	id := c.Param("id")

	log.Printf("[DeleteConnection] Deleting connection: %s", id)

	var connection models.Connection
	if err := c.DB().GetForPublicID(c.AppID(), id, &connection); err != nil {
		log.Printf("[DeleteConnection] Connection not found: %s", id)
		c.NotFound("Connection not found")
		return
	}

	// Clean up the webhook endpoint from the payment provider
	if connection.ExternalWebhookEndpointID != nil && *connection.ExternalWebhookEndpointID != "" {
		log.Printf("[DeleteConnection] Deleting webhook endpoint: %s", *connection.ExternalWebhookEndpointID)

		payment, err := paymentx.New(&connection, c.Crypto())
		if err != nil {
			// Log the error but continue with connection deletion
			log.Printf("[DeleteConnection] Failed to initialize payment client for webhook cleanup: %v", err)
		} else {
			if err := payment.DeleteWebhook(*connection.ExternalWebhookEndpointID); err != nil {
				// Log the error but continue with connection deletion
				// The webhook endpoint may have been manually deleted already
				log.Printf("[DeleteConnection] Failed to delete webhook endpoint (may already be deleted): %v", err)
			} else {
				log.Printf("[DeleteConnection] Webhook endpoint deleted successfully")
			}
		}
	}

	if err := c.DB().DeleteForID(&connection); err != nil {
		log.Printf("[DeleteConnection] Failed to delete connection from database: %v", err)
		c.ServerError("Failed to delete connection", err)
		return
	}

	log.Printf("[DeleteConnection] Connection deleted successfully: %s", id)
	c.OK(map[string]any{"deleted": true, "id": id})
}
