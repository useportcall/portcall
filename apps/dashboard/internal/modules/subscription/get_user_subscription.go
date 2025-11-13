package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetUserSubscription(c *routerx.Context) {
	userID := c.Param("id")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.ServerError("Failed to get user", err)
		return
	}

	var subscription models.Subscription
	if err := c.DB().FindFirst(&subscription, "app_id = ? AND user_id = ? AND status = 'active'", c.AppID(), user.ID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.OK(nil)
			return
		}

		c.ServerError("Failed to get user subscription", err)
		return
	}

	c.OK(new(Subscription).Set(&subscription))
}
