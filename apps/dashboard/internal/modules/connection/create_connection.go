package connection

import (
	"log"

	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateConnectionRequest struct {
	Name      string `json:"name"`
	Source    string `json:"source"`
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
}

func CreateConnection(c *routerx.Context) {
	log.Printf("[CreateConnection] Creating new connection for app %d", c.AppID())

	var body CreateConnectionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[CreateConnection] Invalid request body: %v", err)
		c.BadRequest("Invalid request body")
		return
	}

	log.Printf("[CreateConnection] Connection details: Name=%s, Source=%s",
		body.Name, body.Source)

	encryptedSecretKey, err := c.Crypto().Encrypt(body.SecretKey)
	if err != nil {
		log.Printf("[CreateConnection] Error encrypting secret key: %v", err)
		c.ServerError("Internal server error", err)
		return
	}

	// Generate a public ID for the connection first (needed for webhook URL)
	connectionPublicID := utils.GenPublicID("connect")

	connection := models.Connection{
		PublicID:     connectionPublicID,
		Name:         body.Name,
		Source:       body.Source,
		PublicKey:    body.PublicKey,
		EncryptedKey: encryptedSecretKey,
		AppID:        c.AppID(),
	}

	log.Printf("[CreateConnection] Verifying connection with payment provider...")

	payment, err := paymentx.New(&connection, c.Crypto())
	if err != nil {
		log.Printf("[CreateConnection] Failed to initialize payment client: %v", err)
		c.ServerError("Failed to initialize connection", err)
		return
	}

	if err := payment.Verify(); err != nil {
		log.Printf("[CreateConnection] Failed to verify connection with provider: %v", err)
		c.BadRequest("Failed to verify connection")
		return
	}

	log.Printf("[CreateConnection] Connection verified successfully")

	// Auto-register when supported. Braintree currently uses manual webhook setup.
	webhookURL := buildWebhookURL(body.Source, connectionPublicID)
	log.Printf("[CreateConnection] Registering webhook endpoint: %s", webhookURL)

	webhookReg, err := payment.RegisterWebhook(webhookURL)
	if err != nil {
		log.Printf("[CreateConnection] Failed to register webhook with provider: %v", err)
		c.ServerError("Failed to register webhook with payment provider", err)
		return
	}

	if webhookReg != nil && webhookReg.EndpointID != "" {
		log.Printf("[CreateConnection] Webhook registered successfully: EndpointID=%s", webhookReg.EndpointID)
	}

	if err := attachWebhookRegistration(c, &connection, payment, webhookReg); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if err := c.DB().Create(&connection); err != nil {
		log.Printf("[CreateConnection] Error creating connection in database: %v", err)
		cleanupRegisteredWebhook(payment, webhookReg, "database error")
		c.ServerError("Failed to create connection", err)
		return
	}

	log.Printf("[CreateConnection] Connection created: ID=%d, PublicID=%s, Source=%s",
		connection.ID, connection.PublicID, connection.Source)

	c.OK(connectionResponse(&connection, webhookURL))
}
