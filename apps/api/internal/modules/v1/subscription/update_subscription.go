package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateSubscriptionPayload struct {
	PlanID           string `json:"plan_id"`
	IsFree           *bool  `json:"is_free"`
	ApplyAtNextReset *bool  `json:"apply_at_next_reset"` // If true, schedule plan change for next reset instead of immediate
}

func UpdateSubscription(c *routerx.Context) {
	subscriptionID := c.Param("subscription_id")

	var p UpdateSubscriptionPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.BadRequest("Invalid request payload")
		return
	}

	var subscription models.Subscription
	if err := c.DB().GetForPublicID(c.AppID(), subscriptionID, &subscription); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if subscription.Status != "active" {
		c.BadRequest("Subscription is not active")
		return
	}

	// Handle is_free update (can be done independently of plan change)
	if p.IsFree != nil {
		subscription.IsFree = *p.IsFree
	}

	// If no plan change, just save the is_free update and return
	if p.PlanID == "" {
		if err := c.DB().Save(&subscription); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
		c.OK(buildSubscriptionResponse(c, &subscription))
		return
	}

	var oldPlan models.Plan
	if err := c.DB().FindForID(*subscription.PlanID, &oldPlan); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	var newPlan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), p.PlanID, &newPlan); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if newPlan.Currency != subscription.Currency {
		c.BadRequest("Cannot change currency when switching plan")
		return
	}

	if newPlan.Interval != subscription.BillingInterval {
		c.BadRequest("Cannot change interval when switching plan")
		return
	}

	if newPlan.IntervalCount != subscription.BillingIntervalCount {
		c.BadRequest("Cannot change interval count when switching plan")
		return
	}

	// If apply_at_next_reset is true, schedule the change instead of applying immediately
	if p.ApplyAtNextReset != nil && *p.ApplyAtNextReset {
		subscription.ScheduledPlanID = &newPlan.ID
		if err := c.DB().Save(&subscription); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
		c.OK(buildSubscriptionResponse(c, &subscription))
		return
	}

	// Apply plan switch immediately
	if !applyPlanSwitch(c, &subscription, &oldPlan, &newPlan) {
		return
	}

	c.OK(buildSubscriptionResponse(c, &subscription))
}
