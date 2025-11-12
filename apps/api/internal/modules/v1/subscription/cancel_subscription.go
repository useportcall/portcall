package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CancelSubscription(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")
	resetEntitlements := true

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.ServerError("Internal server error")
		return
	}

	if subscription.Status != "active" {
		c.ServerError("Internal server error")
		return
	}

	subscription.Status = "canceled"
	if err := c.DB().Save(&subscription); err != nil {
		c.ServerError("Internal server error")
		return
	}

	// send notification
	payload := map[string]any{"subscription_id": subscription.ID}
	if err := c.Queue().Enqueue("subscription_canceled", &payload, "email_queue"); err != nil {
		c.ServerError("Internal server error")
		return
	}

	// TODO remove entitlements if setting says so
	if resetEntitlements {
		payload := map[string]any{"subscription_id": subscription.ID}
		if err := c.Queue().Enqueue("start_subscription_reset", &payload, "billing_queue"); err != nil {
			c.ServerError("Internal server error")
			return
		}
	}

	c.OK(new(apix.Subscription).Set(&subscription))
}
