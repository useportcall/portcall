package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func persistSubscriptionWithItems(
	db dbx.IORM,
	sub *models.Subscription,
	plan *models.Plan,
	planItemIDs []uint,
) (int, error) {
	itemCount := 0
	err := db.Txn(func(tx dbx.IORM) error {
		if err := tx.Create(sub); err != nil {
			return err
		}
		if len(planItemIDs) == 0 {
			return nil
		}
		items, err := buildSubItems(tx, plan, sub.ID, planItemIDs)
		if err != nil {
			return err
		}
		for _, item := range items {
			if err := tx.Create(item); err != nil {
				return err
			}
		}
		itemCount = len(items)
		return nil
	})
	return itemCount, err
}

func setRollbackPlan(db dbx.IORM, appID uint, sub *models.Subscription) {
	id, err := findRollbackPlanID(db, appID, sub.PlanID)
	if err == nil && id != nil {
		sub.RollbackPlanID = id
	}
}
