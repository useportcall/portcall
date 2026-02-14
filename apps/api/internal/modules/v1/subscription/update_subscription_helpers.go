package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// applyPlanSwitch replaces subscription items with those from the new plan
// and enqueues the plan-switch billing job.
func applyPlanSwitch(c *routerx.Context, subscription *models.Subscription, oldPlan, newPlan *models.Plan) bool {
	subscription.ScheduledPlanID = nil
	subscription.PlanID = &newPlan.ID
	subscription.InvoiceDueByDays = newPlan.InvoiceDueByDays

	var si models.SubscriptionItem
	if err := c.DB().Delete(&si, "subscription_id = ?", subscription.ID); err != nil {
		c.ServerError("Internal server error", err)
		return false
	}

	if !createSubscriptionItems(c, subscription.ID, newPlan) {
		return false
	}

	if err := c.DB().Save(subscription); err != nil {
		c.ServerError("Internal server error", err)
		return false
	}

	if err := c.Queue().Enqueue(
		"process_plan_switch",
		map[string]any{"old_plan_id": oldPlan.ID, "new_plan_id": newPlan.ID, "subscription_id": subscription.ID},
		"billing_queue",
	); err != nil {
		c.ServerError("Internal server error", err)
		return false
	}

	return true
}

// createSubscriptionItems creates subscription items from plan items.
func createSubscriptionItems(c *routerx.Context, subscriptionID uint, plan *models.Plan) bool {
	var planItems []models.PlanItem
	if err := c.DB().List(&planItems, "plan_id = ?", plan.ID); err != nil {
		c.ServerError("Internal server error", err)
		return false
	}

	for _, pi := range planItems {
		title := pi.PublicTitle
		if pi.PricingModel == "fixed" {
			title = plan.Name
		}

		item := models.SubscriptionItem{
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
			SubscriptionID: subscriptionID,
		}
		if err := c.DB().Create(&item); err != nil {
			c.ServerError("Internal server error", err)
			return false
		}
	}

	return true
}
