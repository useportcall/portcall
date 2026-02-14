package connection

import (
	"fmt"
	"log"
	"os"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func buildWebhookURL(source string, connectionPublicID string) string {
	base := os.Getenv("WEBHOOK_BASE_URL")
	if base == "" {
		base = "https://webhook.useportcall.com"
	}
	path := "stripe"
	if source == "braintree" {
		path = "braintree"
	}
	return fmt.Sprintf("%s/%s/%s", base, path, connectionPublicID)
}

func attachWebhookRegistration(c *routerx.Context, connection *models.Connection, payment paymentx.IPaymentClient, webhookReg *paymentx.WebhookRegistration) error {
	if webhookReg == nil {
		return nil
	}
	if webhookReg.Secret != "" {
		encryptedSecret, err := c.Crypto().Encrypt(webhookReg.Secret)
		if err != nil {
			log.Printf("[CreateConnection] Error encrypting webhook secret: %v", err)
			cleanupRegisteredWebhook(payment, webhookReg, "encryption error")
			return err
		}
		connection.EncryptedWebhookSecret = &encryptedSecret
	}
	if webhookReg.EndpointID != "" {
		connection.ExternalWebhookEndpointID = &webhookReg.EndpointID
	}
	return nil
}

func cleanupRegisteredWebhook(payment paymentx.IPaymentClient, webhookReg *paymentx.WebhookRegistration, reason string) {
	if webhookReg == nil || webhookReg.EndpointID == "" {
		return
	}
	if err := payment.DeleteWebhook(webhookReg.EndpointID); err != nil {
		log.Printf("[CreateConnection] Failed to cleanup webhook after %s: %v", reason, err)
	}
}

func connectionResponse(connection *models.Connection, webhookURL string) *apix.Connection {
	response := new(apix.Connection).Set(connection)
	if connection.Source == "braintree" {
		response.WebhookURL = &webhookURL
	}
	return response
}
