package secret

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateSecret(c *routerx.Context) {
	var body CreateSecretRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.Value == "" {
		c.BadRequest("Missing or invalid value for secret")
		return
	}

	hash, err := c.Crypto().Encrypt(body.Value)
	if err != nil {
		c.ServerError("Failed to encrypt secret value")
		return
	}

	secret := &models.Secret{
		PublicID: dbx.GenPublicID("sk"),
		AppID:    c.AppID(),
		KeyHash:  hash,
	}

	if err := c.DB().Create(secret); err != nil {
		c.ServerError("Failed to create secret")
		return
	}

	c.OK(new(Secret).Set(secret))
}

type CreateSecretRequest struct {
	Value string `json:"value"`
}
