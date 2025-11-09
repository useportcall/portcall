package address

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteAddress(c *routerx.Context) {
	id := c.Param("id")

	address := &models.Address{}
	if err := c.DB().GetForPublicID(c.AppID(), id, address); err != nil {
		c.NotFound("Address not found")
		return
	}

	if err := c.DB().RemoveForPublicID(c.AppID(), id, address); err != nil {
		c.ServerError("Failed to delete address")
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
