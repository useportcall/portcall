package user

import (
	"log"

	"github.com/useportcall/portcall/apps/dashboard/portcall"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteUser(c *routerx.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), id, &user); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("User not found")
			return
		}
		c.ServerError("Failed to get user", err)
		return
	}

	if err := c.DB().RemoveForPublicID(c.AppID(), id, &user); err != nil {
		c.ServerError("Failed to delete user", err)
		return
	}

	// decrement df number of users by 1
	if err := c.Queue().Enqueue("df_decrement", map[string]any{"user_id": c.PublicAppID(), "feature": portcall.Features.NumberOfUsers, "is_test": !c.IsLive()}, "billing_queue"); err != nil {
		log.Printf("Error enqueueing df_decrement: %v", err)
		c.ServerError("error updating feature usage", err)
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
