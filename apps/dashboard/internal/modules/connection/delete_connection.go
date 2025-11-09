package connection

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteConnection(c *routerx.Context) {
	id := c.Param("id")

	var connection models.Connection
	if err := c.DB().GetForPublicID(c.AppID(), id, &connection); err != nil {
		c.NotFound("Connection not found")
		return
	}

	if err := c.DB().DeleteForID(&connection); err != nil {
		c.ServerError("Failed to delete connection")
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
