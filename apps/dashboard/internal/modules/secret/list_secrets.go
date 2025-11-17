package secret

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListSecrets(c *routerx.Context) {
	secrets := []models.Secret{}
	if err := c.DB().ListForAppID(c.AppID(), &secrets, nil); err != nil {
		c.ServerError("Failed to list secrets", err)
		return
	}

	response := make([]apix.Secret, len(secrets))
	for i, secret := range secrets {
		response[i].Set(&secret)
	}

	c.OK(response)
}
