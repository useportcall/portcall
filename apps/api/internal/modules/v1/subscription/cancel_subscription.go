package subscription

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CancelSubscriptionRequest struct {
	SkipEntitlementReset bool `json:"skip_entitlement_reset"`
}

func CancelSubscription(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")

	var body CancelSubscriptionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if subscription.Status != "active" {
		c.ServerError("Internal server error", fmt.Errorf("only active subscriptions can be canceled, status: %s", subscription.Status))
		return
	}

	subscription.Status = "canceled"
	if err := c.DB().Save(&subscription); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	payload := map[string]any{"subscription_id": subscription.ID}
	if err := c.Queue().Enqueue("subscription_canceled", &payload, "email_queue"); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if body.SkipEntitlementReset {
		log.Printf("skipping entitlement reset for subscription %s as per request", subscription.PublicID)
		c.OK(new(apix.Subscription).Set(&subscription))
		return
	}

	payload = map[string]any{"subscription_id": subscription.ID}
	if err := c.Queue().Enqueue("start_subscription_reset", &payload, "billing_queue"); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	c.OK(new(apix.Subscription).Set(&subscription))
}
