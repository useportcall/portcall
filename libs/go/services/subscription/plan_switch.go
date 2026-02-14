package subscription

import (
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// PlanSwitch determines whether a plan change is an upgrade or
// downgrade and calculates the prorated price difference for upgrades.
func (s *service) PlanSwitch(input *PlanSwitchInput) (*PlanSwitchResult, error) {
	log.Printf("Processing plan switch old=%d new=%d sub=%d",
		input.OldPlanID, input.NewPlanID, input.SubscriptionID)

	oldItem, err := findFixedPlanItem(s.db, input.OldPlanID)
	if err != nil {
		return nil, err
	}

	newItem, err := findFixedPlanItem(s.db, input.NewPlanID)
	if err != nil {
		return nil, err
	}

	if newItem.UnitAmount <= oldItem.UnitAmount {
		return &PlanSwitchResult{
			IsUpgrade:      false,
			SubscriptionID: input.SubscriptionID,
		}, nil
	}

	var sub models.Subscription
	if err := s.db.FindForID(input.SubscriptionID, &sub); err != nil {
		return nil, err
	}

	diff := newItem.UnitAmount - oldItem.UnitAmount
	prorated := Prorate(diff, sub.LastResetAt, sub.NextResetAt, time.Now())

	return &PlanSwitchResult{
		IsUpgrade:       true,
		SubscriptionID:  input.SubscriptionID,
		UserID:          sub.UserID,
		OldPlanID:       input.OldPlanID,
		NewPlanID:       input.NewPlanID,
		PriceDifference: prorated,
	}, nil
}

func findFixedPlanItem(db interface{ FindFirst(any, ...any) error },
	planID uint,
) (*models.PlanItem, error) {
	var item models.PlanItem
	if err := db.FindFirst(&item, "plan_id = ? AND pricing_model = ?",
		planID, "fixed"); err != nil {
		return nil, err
	}
	return &item, nil
}
