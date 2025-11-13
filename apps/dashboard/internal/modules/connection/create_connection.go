package connection

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateConnection(c *routerx.Context) {
	var body CreateConnectionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	encryptedSecretKey, err := c.Crypto().Encrypt(body.SecretKey)
	if err != nil {
		c.ServerError("Internal server error", err)
	}

	encryptedWebhookSecret, err := c.Crypto().Encrypt(body.WebhookSecret)
	if err != nil {
		c.ServerError("Internal server error", err)
	}

	connection := models.Connection{
		PublicID:               utils.GenPublicID("connect"),
		Name:                   body.Name,
		Source:                 body.Source,
		PublicKey:              body.PublicKey,
		EncryptedKey:           encryptedSecretKey,
		EncryptedWebhookSecret: &encryptedWebhookSecret,
		AppID:                  c.AppID(),
	}

	payment, err := paymentx.New(&connection, c.Crypto())
	if err != nil {
		c.ServerError("Failed to initialize connection", err)
		return
	}

	if err := payment.Verify(); err != nil {
		c.BadRequest("Failed to verify connection")
	}

	if err := c.DB().Create(&connection); err != nil {
		c.ServerError("Failed to create connection", err)
		return
	}

	c.OK(new(Connection).Set(&connection))
}
