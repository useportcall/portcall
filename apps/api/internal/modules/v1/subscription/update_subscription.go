package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateSubscriptionPayload struct {
	PlanID string `json:"plan_id"`
}

func UpdateSubscription(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")

	var p UpdateSubscriptionPayload
	if err := c.BindJSON(&p); err != nil {
		c.BadRequest("Invalid request payload")
		return
	}

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.ServerError("Internal server error")
		return
	}

	if subscription.Status != "active" {
		c.BadRequest("Subscription is not active")
		return
	}

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), p.PlanID, &plan); err != nil {
		c.ServerError("Internal server error")
		return
	}

	if plan.Currency != subscription.Currency {
		c.BadRequest("Cannot change currency when switching plan")
		return
	}

	if plan.Interval != subscription.BillingInterval {
		c.BadRequest("Cannot change interval when switching plan")
		return
	}

	if plan.IntervalCount != subscription.BillingIntervalCount {
		c.BadRequest("Cannot change interval count when switching plan")
		return
	}

	subscription.PlanID = &plan.ID
	subscription.InvoiceDueByDays = plan.InvoiceDueByDays

	var si models.SubscriptionItem
	if err := c.DB().Delete(&si, "subscription_id = ?", subscription.ID); err != nil {
		c.ServerError("Internal server error")
		return
	}

	if err := c.DB().Save(&subscription); err != nil {
		c.ServerError("Internal server error")
		return
	}

	payload := map[string]any{"user_id": subscription.UserID, "plan_id": plan.ID}
	if err := c.Queue().Enqueue("create_entitlements", payload, "billing_queue"); err != nil {
		c.ServerError("Internal server error")
		return
	}

	c.OK(new(apix.Subscription).Set(&subscription))
}
