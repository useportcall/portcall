package user

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteUser(c *routerx.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), id, &user); err != nil {
		c.ServerError("Failed to get user", err)
	}

	if err := c.DB().RemoveForPublicID(c.AppID(), id, &user); err != nil {
		c.ServerError("Failed to delete user", err)
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
