package subscription

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
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
		c.ServerError("Internal server error", err)
		return
	}

	if subscription.Status != "active" {
		c.BadRequest("Subscription is not active")
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

	subscription.PlanID = &newPlan.ID
	subscription.InvoiceDueByDays = newPlan.InvoiceDueByDays

	var si models.SubscriptionItem
	if err := c.DB().Delete(&si, "subscription_id = ?", subscription.ID); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", newPlan.ID); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	for _, pi := range planItems {
		var title string
		if pi.PricingModel == "fixed" {
			title = newPlan.Name
		} else {
			title = pi.PublicTitle
		}

		subscriptionItem := models.SubscriptionItem{
			PublicID:       dbx.GenPublicID("si"),
			PlanItemID:     &pi.ID,
			Quantity:       pi.Quantity,
			AppID:          pi.AppID,
			UnitAmount:     pi.UnitAmount,
			PricingModel:   pi.PricingModel,
			Tiers:          pi.Tiers,
			Maximum:        pi.Maximum,
			Minimum:        pi.Minimum,
			Title:          title,
			Description:    pi.PublicDescription,
			SubscriptionID: subscription.ID,
		}
		if err := c.DB().Create(&subscriptionItem); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
	}

	if err := c.DB().Save(&subscription); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if err := c.Queue().Enqueue(
		"process_plan_switch",
		map[string]any{"old_plan_id": oldPlan.ID, "new_plan_id": newPlan.ID, "subscription_id": subscription.ID},
		"billing_queue",
	); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	if err := c.Queue().Enqueue(
		"create_entitlements",
		map[string]any{"user_id": subscription.UserID, "plan_id": newPlan.ID},
		"billing_queue",
	); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	c.OK(new(apix.Subscription).Set(&subscription))
}
