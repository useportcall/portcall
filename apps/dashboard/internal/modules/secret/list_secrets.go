package secret

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListSecrets(c *routerx.Context) {
	secrets := []models.Secret{}
	if err := c.DB().ListForAppID(c.AppID(), &secrets, nil); err != nil {
		c.ServerError("Failed to list secrets", err)
		return
	}

	response := make([]*Secret, len(secrets))
	for i := range secrets {
		response[i] = new(Secret).Set(&secrets[i])
	}

	c.OK(response)
}
