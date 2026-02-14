package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// buildSubscriptionResponse builds the full subscription response with plan details.
func buildSubscriptionResponse(c *routerx.Context, subscription *models.Subscription) *apix.Subscription {
	result := new(apix.Subscription).Set(subscription)

	if subscription.PlanID != nil {
		var plan models.Plan
		if err := c.DB().FindForID(*subscription.PlanID, &plan); err == nil {
			result.Plan = new(apix.Plan).Set(&plan)
			result.PlanID = plan.PublicID
		}
	}

	if subscription.ScheduledPlanID != nil {
		var scheduledPlan models.Plan
		if err := c.DB().FindForID(*subscription.ScheduledPlanID, &scheduledPlan); err == nil {
			result.ScheduledPlan = new(apix.Plan).Set(&scheduledPlan)
			result.ScheduledPlanID = &scheduledPlan.PublicID
		}
	}

	return result
}
