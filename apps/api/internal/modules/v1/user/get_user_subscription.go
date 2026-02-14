package user

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// UserSubscriptionResponse extends apix.Subscription with has_payment_method
type UserSubscriptionResponse struct {
	*apix.Subscription
	HasPaymentMethod bool `json:"has_payment_method"`
}

// GetUserSubscription returns the current billable subscription for a user.
func GetUserSubscription(c *routerx.Context) {
	userID := c.Param("id")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("User not found")
			return
		}
		c.ServerError("Failed to get user", err)
		return
	}

	var subscription models.Subscription
	if err := c.DB().FindFirst(&subscription,
		"app_id = ? AND user_id = ? AND status IN (?, ?)",
		c.AppID(), user.ID, "active", "past_due",
	); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("No billable subscription found for user")
			return
		}
		c.ServerError("Failed to get subscription", err)
		return
	}

	result := new(apix.Subscription).Set(&subscription)
	result.User = new(apix.User).Set(&user)

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

	var paymentMethod models.PaymentMethod
	hasPaymentMethod := false
	if err := c.DB().FindFirst(&paymentMethod, "app_id = ? AND user_id = ?", c.AppID(), user.ID); err == nil {
		hasPaymentMethod = true
	}

	c.OK(UserSubscriptionResponse{
		Subscription:     result,
		HasPaymentMethod: hasPaymentMethod,
	})
}
