package subscription

import "github.com/useportcall/portcall/libs/go/dbx/models"

func cheapestFixedPlan(db interface {
	FindForID(uint, any) error
	FindFirst(any, ...any) error
}, ids []uint, current *uint) (*uint, error) {
	var (
		bestID    *uint
		bestPrice int64
	)
	for _, id := range ids {
		if current != nil && id == *current {
			continue
		}
		var plan models.Plan
		if err := db.FindForID(id, &plan); err != nil {
			continue
		}
		var item models.PlanItem
		if err := db.FindFirst(&item, "plan_id = ? AND pricing_model = ?", plan.ID, "fixed"); err != nil {
			continue
		}
		if bestID == nil || item.UnitAmount < bestPrice {
			candidate := plan.ID
			bestID, bestPrice = &candidate, item.UnitAmount
		}
	}
	return bestID, nil
}
