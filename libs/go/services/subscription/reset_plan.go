package subscription

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func applyScheduledPlan(db dbx.IORM, sub *models.Subscription) (*uint, error) {
	if sub.ScheduledPlanID == nil {
		return nil, nil
	}

	var nextPlan models.Plan
	if err := db.FindForID(*sub.ScheduledPlanID, &nextPlan); err != nil {
		return nil, err
	}
	if err := validatePlanSwitch(sub, &nextPlan); err != nil {
		return nil, err
	}

	var si models.SubscriptionItem
	if err := db.Delete(&si, "subscription_id = ?", sub.ID); err != nil {
		return nil, err
	}
	if err := createItemsForPlan(db, sub.ID, &nextPlan); err != nil {
		return nil, err
	}

	sub.PlanID = &nextPlan.ID
	sub.ScheduledPlanID = nil
	sub.InvoiceDueByDays = nextPlan.InvoiceDueByDays
	sub.BillingInterval = nextPlan.Interval
	sub.BillingIntervalCount = nextPlan.IntervalCount
	sub.IsFree = nextPlan.IsFree

	id := nextPlan.ID
	return &id, nil
}

func createItemsForPlan(db dbx.IORM, subscriptionID uint, plan *models.Plan) error {
	var planItemIDs []uint
	if err := db.ListIDs("plan_items", &planItemIDs, "plan_id = ?", plan.ID); err != nil {
		return err
	}
	if len(planItemIDs) == 0 {
		return nil
	}

	items, err := buildSubItems(db, plan, subscriptionID, planItemIDs)
	if err != nil {
		return fmt.Errorf("build reset items: %w", err)
	}
	for _, item := range items {
		if err := db.Create(item); err != nil {
			return err
		}
	}
	return nil
}
