package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetSubscription(c *routerx.Context) {
	subscriptionID := c.Param("id")

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("Subscription not found")
			return
		}

		c.ServerError("Failed to get subscription", err)
		return
	}

	result := new(apix.Subscription).Set(&subscription)

	var user models.User
	if err := c.DB().FindForID(subscription.UserID, &user); err != nil {
		c.ServerError("Failed to get user for subscription", err)
		return
	}

	result.User = new(apix.User).Set(&user)

	if subscription.PlanID != nil {
		var plan models.Plan
		if err := c.DB().FindForID(*subscription.PlanID, &plan); err != nil {
			c.ServerError("Failed to get plan for subscription", err)
			return
		}

		result.Plan = new(apix.Plan).Set(&plan)
	}

	c.OK(result)
}
