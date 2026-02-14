package subscription

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// CreateItems bulk-creates all subscription items for a plan in a
// single transaction. Replaces the old start→iterate pattern.
func (s *service) CreateItems(input *CreateItemsInput) (*CreateItemsResult, error) {
	log.Printf("Creating items for subscription %d plan %d",
		input.SubscriptionID, input.PlanID)

	planItemIDs, err := listPlanItemIDs(s.db, input.PlanID)
	if err != nil {
		return nil, err
	}
	if len(planItemIDs) == 0 {
		log.Printf("No plan items for plan %d", input.PlanID)
		return &CreateItemsResult{
			SubscriptionID: input.SubscriptionID,
		}, nil
	}

	plan, err := findPlan(s.db, input.PlanID)
	if err != nil {
		return nil, err
	}

	items, err := buildSubItems(s.db, plan, input.SubscriptionID, planItemIDs)
	if err != nil {
		return nil, err
	}

	if err := persistItems(s.db, items); err != nil {
		return nil, err
	}

	log.Printf("Created %d items for subscription %d", len(items), input.SubscriptionID)
	return &CreateItemsResult{
		SubscriptionID: input.SubscriptionID,
		ItemCount:      len(items),
	}, nil
}

func buildSubItems(
	db dbx.IORM, plan *models.Plan,
	subID uint, planItemIDs []uint,
) ([]*models.SubscriptionItem, error) {
	items := make([]*models.SubscriptionItem, 0, len(planItemIDs))
	for _, id := range planItemIDs {
		var pi models.PlanItem
		if err := db.FindForID(id, &pi); err != nil {
			return nil, fmt.Errorf("plan item %d: %w", id, err)
		}
		title := pi.PublicTitle
		if pi.PricingModel == "fixed" {
			title = plan.Name
		}
		items = append(items, &models.SubscriptionItem{
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
			SubscriptionID: subID,
			Interval:       pi.Interval,
			IntervalCount:  pi.IntervalCount,
		})
	}
	return items, nil
}

func listPlanItemIDs(db dbx.IORM, planID uint) ([]uint, error) {
	var ids []uint
	if err := db.ListIDs("plan_items", &ids, "plan_id = ?", planID); err != nil {
		return nil, err
	}
	return ids, nil
}

func persistItems(db dbx.IORM, items []*models.SubscriptionItem) error {
	return db.Txn(func(tx dbx.IORM) error {
		for _, item := range items {
			if err := tx.Create(item); err != nil {
				return fmt.Errorf("create item %s: %w", item.PublicID, err)
			}
		}
		return nil
	})
}
