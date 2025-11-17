package secret

import (
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DisableSecret(c *routerx.Context) {
	id := c.Param("id")
	if id == "" {
		c.BadRequest("Missing secret_id parameter")
		return
	}

	secret := &models.Secret{}
	if err := c.DB().GetForPublicID(c.AppID(), id, secret); err != nil {
		c.NotFound("Secret not found")
		return
	}

	now := time.Now()
	secret.DisabledAt = &now

	if err := c.DB().UpdateForPublicID(c.AppID(), secret.PublicID, secret); err != nil {
		c.ServerError("Failed to disable secret", err)
		return
	}

	c.OK(new(apix.Secret).Set(secret))
}
