package subscription

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/apps/dashboard/portcall"
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
		log.Printf("Invalid request body: %v", err)
		c.BadRequest("invalid request body")
		return
	}

	log.Printf("Canceling subscription %s", subscriptionID)

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

	// decrement df max subscriptions by 1
	if err := c.Queue().Enqueue("df_decrement", map[string]any{"user_id": c.PublicAppID(), "feature": portcall.Features.MaxSubscriptions, "is_test": !c.IsLive()}, "billing_queue"); err != nil {
		log.Printf("Error enqueueing df_decrement: %v", err)
		c.ServerError("error updating feature usage", err)
		return
	}

	c.OK(new(apix.Subscription).Set(&subscription))
}
