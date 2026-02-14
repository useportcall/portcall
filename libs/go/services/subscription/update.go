package subscription

import (
	"errors"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Update switches a subscription to a new plan. It validates plan
// compatibility, deletes old items, and updates the subscription
// in a single transaction.
func (s *service) Update(input *UpdateInput) (*UpdateResult, error) {
	log.Printf("Updating subscription %d to plan %d",
		input.SubscriptionID, input.PlanID)

	var sub models.Subscription
	if err := s.db.FindForID(input.SubscriptionID, &sub); err != nil {
		return nil, err
	}

	oldPlan, err := findPlan(s.db, *sub.PlanID)
	if err != nil {
		return nil, err
	}

	newPlan, err := findPlan(s.db, input.PlanID)
	if err != nil {
		return nil, err
	}

	if err := validatePlanSwitch(&sub, newPlan); err != nil {
		return nil, err
	}

	if err := s.applyPlanUpdate(&sub, newPlan); err != nil {
		return nil, err
	}

	return &UpdateResult{
		Subscription: &sub,
		OldPlanID:    oldPlan.ID,
		NewPlanID:    newPlan.ID,
	}, nil
}

func validatePlanSwitch(sub *models.Subscription, p *models.Plan) error {
	if p.Currency != sub.Currency {
		return errors.New("cannot change currency when switching plan")
	}
	if p.Interval != sub.BillingInterval {
		return errors.New("cannot change interval when switching plan")
	}
	if p.IntervalCount != sub.BillingIntervalCount {
		return errors.New("cannot change interval count when switching plan")
	}
	return nil
}

func (s *service) applyPlanUpdate(
	sub *models.Subscription, newPlan *models.Plan,
) error {
	return s.db.Txn(func(tx dbx.IORM) error {
		var si models.SubscriptionItem
		if err := tx.Delete(&si, "subscription_id = ?", sub.ID); err != nil {
			return err
		}
		sub.ScheduledPlanID = nil
		sub.PlanID = &newPlan.ID
		sub.InvoiceDueByDays = newPlan.InvoiceDueByDays
		return tx.Save(sub)
	})
}
